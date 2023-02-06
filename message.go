package paxos

type messageType int

const (
	Prepare messageType = iota + 1
	Promise
	Propose
	Accept
)

type message struct {
	from   int
	to     int
	mtype  messageType
	seq    int
	preSeq int
	value  string
}

func (m *message) getProposeValue() string {
	return m.value
}

func (m *message) getProposeSeq() int {
	return m.seq
}
