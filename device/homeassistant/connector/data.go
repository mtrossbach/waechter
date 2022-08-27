package connector

type Data struct {
	Id  uint64
	Msg []byte
	Err error
}

type SetId interface {
	SetId(seq uint64)
}
