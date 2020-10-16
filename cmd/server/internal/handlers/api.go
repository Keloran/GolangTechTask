package handlers

type CreateBuff struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

type Buff struct {
	ID       uint64   `json:"id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

type CreateStream struct {
	StreamName string `json:"stream_name"`
}

type StreamAttachBuffs struct {
  Buffs []uint64 `json:"buffs"`
}

type StreamList struct {
  ID uint64 `json:"id"`
  StreamName string `json:"stream_name"`
}

type Stream struct {
  ID uint64 `json:"id"`
  StreamName string `json:"stream_name"`
  Buffs []Buff `json:"buffs"`
}
