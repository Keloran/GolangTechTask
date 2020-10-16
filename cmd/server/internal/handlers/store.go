package handlers

type BuffStore interface {
	GetBuff(id uint64) (*Buff, error)
	SetBuff(*Buff) (id uint64, err error)
	ListBuffs() (map[uint64]*Buff, error)
}

type Store interface {
	BuffStore
	StreamStore
}

type StreamStore interface {
  GetStream(id uint64) (*Stream, error)
  SetStream(*Stream) (id uint64, err error)
  ListStreams() ([]*StreamList, error)
  AppendBuffs(id uint64, buffs []uint64) (*Stream, error)
}
