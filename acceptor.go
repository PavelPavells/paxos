package paxos

import "log"

type acceptor struct {
	id             int
	learners       []int
	acceptMessage  message
	promiseMessage message
	nodeNetwork    nodeNetwork
}

func NewAcceptor(id int, nodeNetwork nodeNetwork, learners ...int) acceptor {
	newAcceptor := acceptor{id: id, nodeNetwork: nodeNetwork}
	newAcceptor.learners = learners

	return newAcceptor
}

func (a *acceptor) run() {
	for {
		m := a.nodeNetwork.receive()

		if m == nil {
			continue
		}

		switch m.mtype {
		case Prepare:
			promiseMessage := a.receivePrepare(*m)
			a.nodeNetwork.send(*promiseMessage)

			continue
		case Propose:
			accepted := a.receivePropose(*m)

			if accepted {
				for _, lId := range a.learners {
					m.from = a.id
					m.to = lId
					m.mtype = Accept
					a.nodeNetwork.send(*m)
				}
			}
		default:
			log.Fatalln("Unsupport message in acceptor ID: ", a.id)
		}
	}

	log.Println("accetor :", a.id, " leave.")
}

func (a *acceptor) receivePrepare(prepare message) *message {
	if a.promiseMessage.getProposeSeq() >= prepare.getProposeSeq() {
		log.Println("ID:", a.id, "Already accept bigger one")

		return nil
	}

	log.Println("ID: ", a.id, " Promise")

	prepare.to = prepare.from
	prepare.from = a.id
	prepare.mtype = Promise
	a.acceptMessage = prepare

	return &prepare
}

func (a *acceptor) receivePropose(proposeMessage message) bool {
	log.Println("accept:check propose. ", a.acceptMessage.getProposeSeq(), proposeMessage.getProposeSeq())

	if a.acceptMessage.getProposeSeq() > proposeMessage.getProposeSeq() || a.acceptMessage.getProposeSeq() < proposeMessage.getProposeSeq() {
		log.Println("ID:", a.id, " acceptor not take propose:", proposeMessage.value)

		return false
	}

	return true
}
