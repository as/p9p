package p9p

type Fid uint32

type Qid struct {
	Type byte
	Ver  uint32
	Path uint64
}

type Qids struct {
	List []Qid
}
