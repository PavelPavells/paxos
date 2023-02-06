package paxos

import "log"

type learner struct {
	id               int
	acceptedMessages map[int]message
	nodeNetwork      nodeNetwork
}

func NewLearner(id int, nodeNetwork nodeNetwork, acceptorIDs ...int) *learner {
	newLearner := &learner{id: id, nodeNetwork: nodeNetwork}
	newLearner.acceptedMessages = make(map[int]message)

	for _, acceptID := range acceptorIDs {
		newLearner.acceptedMessages[acceptID] = message{}
	}

	return newLearner
}

func (learner *learner) run() string {
	for {
		message := learner.nodeNetwork.receive()

		if message == nil {
			continue
		}

		log.Println("Learner: receive message: ", *message)
		learner.handleReceiveAccept(*message)
		learnMessage, hasLearn := learner.chosen()

		if !hasLearn {
			continue
		}

		return learnMessage.getProposeValue()
	}
}

func (learner *learner) handleReceiveAccept(acceptMessage message) {
	hasAcceptedMessages := learner.acceptedMessages[acceptMessage.from]

	if hasAcceptedMessages.getProposeSeq() < acceptMessage.getProposeSeq() {
		learner.acceptedMessages[acceptMessage.from] = acceptMessage
	}
}

func (learner *learner) chosen() (message, bool) {
	acceptCount := make(map[int]int)
	acceptMessages := make(map[int]message)

	for _, message := range learner.acceptedMessages {
		proposalNum := message.getProposeSeq()
		acceptCount[proposalNum]++
		acceptMessages[proposalNum] = message
	}

	for chosenNum, chosenMessage := range acceptMessages {
		if acceptCount[chosenNum] > learner.majority() {
			return chosenMessage, true
		}
	}

	return message{}, false
}

func (learner *learner) majority() int {
	return len(learner.acceptedMessages) / 2 + 1
}
