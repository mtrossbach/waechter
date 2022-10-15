package connection

type subscription struct {
	handler StateEventHandler
	seqId   uint64
}

type SetId interface {
	SetId(seq uint64)
}
