package paxos

import "log"

type proposer struct {
	id           int
	seq          int
	proposeNumber   int
	proposeValue string
	acceptors    map[int]message
	nodeNetwork  nodeNetwork
}

func NewProposer(id int, value string, nodeNetwork nodeNetwork, acceptors ...int) *proposer {
	pr := proposer{id: id, proposeValue: value, seq: 0, nodeNetwork: nodeNetwork}
	pr.acceptors = make(map[int]message, len(acceptors))

	log.Println("Proposer has ", len(acceptors), " acceptors, value: ", pr.proposeValue)

	for _, acceptor := range acceptors {
		pr.acceptors[acceptor] = message{}
	}

	return &pr
}

func (proposer *proposer) run() {
	log.Println("Proposer start run... value:", proposer.proposeValue)

	for !proposer.majorityReached() {
		log.Println("[Proposer:Prepare]")

		outMessages := proposer.prepare()
		log.Println("[Proposer: prepare ", len(outMessages), "message")

		for _, message := range outMessages {
			proposer.nodeNetwork.send(message)
			log.Println("[Proposer: send", message)
		}

		log.Println("[Proposer: prepare recieve..")

		message := proposer.nodeNetwork.receive()

		if message == nil {
			log.Println("[Proposer: no msg... ")

			continue
		}

		log.Println("[Proposer: recev", message)

		switch message.mtype {
		case Promise:
			log.Println(" proposer recev a promise from ", message.from)
			proposer.checkReceivePromise(*message)
		default:
			panic("Unsupport message.")
		}
	}

	log.Println("[Proposer:Propose]")

	log.Println("Proposor propose seq:", proposer.getProposeNum(), " value:", proposer.proposeValue)
	proposeMessages := proposer.propose()

	for _, message := range proposeMessages {
		proposer.nodeNetwork.send(message)
	}
}

func (proposer *proposer) propose() []message {
	sendMessageCount := 0
	var messageList []message

	log.Println("proposer: propose msg:", len(proposer.acceptors))

	for acceptId, acceptMessage := range proposer.acceptors {
		log.Println("check promise id:", acceptMessage.getProposeSeq(), proposer.getProposeNum())

		if acceptMessage.getProposeSeq() == proposer.getProposeNum() {
			msg := message{from: proposer.id, to: acceptId, mtype: Propose, seq: proposer.getProposeNum()}
			msg.value = proposer.proposeValue

			log.Println("Propose val:", msg.value)

			messageList = append(messageList, msg)
		}

		sendMessageCount++

		if sendMessageCount > proposer.majority() {
			break
		}
	}

	log.Println(" proposer propose msg list:", messageList)

	return messageList
}

func (proposer *proposer) prepare() []message {
	proposer.seq++

	sendMessageCount := 0
	var msgList []message

	log.Println("proposer: prepare major msg:", len(proposer.acceptors))

	for acceptId, _ := range proposer.acceptors {
		msg := message{from: proposer.id, to: acceptId, mtype: Prepare, seq: proposer.getProposeNum(), value: proposer.proposeValue}
		msgList = append(msgList, msg)
		sendMessageCount++

		if sendMessageCount > proposer.majority() {
			break
		}
	}
	return msgList
}

func (proposer *proposer) checkReceivePromise(promise message) {
	previousPromise := proposer.acceptors[promise.from]

	log.Println(" previousMessage:", previousPromise, " promiseMessage:", promise)
	log.Println(previousPromise.getProposeSeq(), promise.getProposeSeq())

	if previousPromise.getProposeSeq() < promise.getProposeSeq() {

		log.Println("Proposor:", proposer.id, " get new promise:", promise)

		proposer.acceptors[promise.from] = promise

		if promise.getProposeSeq() > proposer.getProposeNum() {
			proposer.proposeNumber = promise.getProposeSeq()
			proposer.proposeValue = promise.getProposeValue()
		}
	}
}

func (proposer *proposer) majority() int {
	return len(proposer.acceptors) / 2 + 1
}

func (proposer *proposer) getRecevPromiseCount() int {
	receiveCount := 0

	for _, acceptMessage := range proposer.acceptors {
		log.Println(" proposer has total ", len(proposer.acceptors), " acceptor ", acceptMessage, " current Number:", proposer.getProposeNum(), " messageNum:", acceptMessage.getProposeSeq())

		if acceptMessage.getProposeSeq() == proposer.getProposeNum() {
			log.Println("receive ++", receiveCount)
			receiveCount++
		}
	}
	log.Println("Current proposer recev promise count = ", receiveCount)

	return receiveCount
}

func (proposer *proposer) majorityReached() bool {
	return proposer.getRecevPromiseCount() > proposer.majority()
}

func (proposer *proposer) getProposeNum() int {
	proposer.proposeNumber = proposer.seq << 4 | proposer.id

	return proposer.proposeNumber
}
