package paxos

import (
	"log"
	"testing"
	"time"
)

func TestBasicNetwork(t *testing.T) {
	log.Println("Test Basic Network ...............")

	network := CreateNetwork(1, 3, 5, 2, 4)

	go func() {
		network.receiveFrom(5)
		network.receiveFrom(1)
		network.receiveFrom(3)
		network.receiveFrom(2)

		message := network.receiveFrom(4)

		if message == nil {
			t.Errorf("No message detected")
		}
	}()

	firstMessage := message{from: 3, to: 1, mtype: Prepare, seq: 1, preSeq: 0, value: "firstMessage"}
	network.sendTo(firstMessage)

	secondMessage := message{from: 5, to: 3, mtype: Accept, seq: 2, preSeq: 1, value: "secondMessage"}
	network.sendTo(secondMessage)

	thirdMessage := message{from: 4, to: 2, mtype: Promise, seq: 3, preSeq: 2, value: "thirdMessage"}
	network.sendTo(thirdMessage)
	
	time.Sleep(time.Second)
}

func TestSingleProser(t *testing.T) {
	log.Println("Test Proser Function ...............")

	network := CreateNetwork(100, 1, 2, 3, 200)

	var acceptors acceptor[]
	aId := 1

	for aId <= 3 {
		acctor := NewAcceptor(aId, network.getNodeNetwork(aId), 200)
		acceptors = append(acceptors, acctor)
		aId++
	}

	proposer := NewProposer(100, "value 1", network.getNodeNetwork(100), 1, 2, 3)

	go proposer.run()

	for index, _ := range acceptors {
		go acceptor[index].run()
	}

	learner := NewLearner(200, network.getNodeNetwork(200), 1, 2, 3)
	learnValue := learner.run()

	if learnValue != "value1" {
		t.Errorf("Learner learn wrong proposal")
	}
}

func TestTwoProsers(t *testing.T) {
	log.Println("Test Proser Function ...............")

	network := CreateNetwork(100, 1, 2, 3, 200, 101)

	var acceptors acceptor[]

	aId := 1

	for aId <= 3 {
		acctor := NewAcceptor(aId, network.getNodeNetwork(aId), 200)
		acceptors = append(acceptors, acctor)
		aId++
	}

	firstProposer := NewProposer(100, "ExpectValue", network.getNodeNetwork(100), 1, 2, 3)
	go firstProposer.run()

	secondProposer := NewProposer(101, "WrongValue", network.getNodeNetwork(101), 1, 2, 3)
	time.AfterFunc(time.Second, func() {
		secondProposer.run()
	})

	for _, index := range acceptors {
		go acceptors[index].run()
	}

	learner := NewLearner(200, network.getNodeNetwork(200), 1, 2, 3)
	learnValue := learner.run()
	
	if learnValue != "ExpectValue" {
		t.Errorf("Learner learn wrong proposal. Expect:'ExpectValue', learnValue: %v", learnValue)
	}
}
