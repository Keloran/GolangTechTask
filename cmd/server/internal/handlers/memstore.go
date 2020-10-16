package handlers

import (
  "errors"
  "sync"
)

type inMemStore struct {
	mu        *sync.RWMutex

	buffCounter uint64
	buffs     map[uint64]*Buff

	streamCounter uint64
	streams   map[uint64]*Stream
}

func (i *inMemStore) GetBuff(id uint64) (*Buff, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	b, ok := i.buffs[id]
	if !ok {
		return nil, errors.New("buff not found")
	}
	return b, nil
}

func (i *inMemStore) SetBuff(b *Buff) (uint64, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.buffCounter++
	b.ID = i.buffCounter
	i.buffs[i.buffCounter] = b
	return i.buffCounter, nil
}

func (i *inMemStore) ListBuffs() (map[uint64]*Buff, error) {
	return i.buffs, nil
}

func (i *inMemStore) GetStream(id uint64) (*Stream, error) {
  i.mu.RLock()
  defer i.mu.RUnlock()

  b, ok := i.streams[id]
  if !ok {
    return nil, errors.New("stream not found")
  }
  return b, nil
}

func (i *inMemStore) SetStream(s *Stream) (uint64, error) {
  i.mu.Lock()
  defer i.mu.Unlock()

  i.streamCounter++
  s.ID = i.streamCounter
  i.streams[i.streamCounter] = s
  return i.streamCounter, nil
}

func (i *inMemStore) AppendBuffs(id uint64, buffs []uint64) (*Stream, error) {
  i.mu.Lock()
  defer i.mu.Unlock()

  stream, ok := i.streams[id]
  if !ok {
    return nil, errors.New("stream not found")
  }
  bs := []Buff{}
  for _, b := range buffs {
    bs = append(bs, *i.buffs[b])
  }

  ns := Stream{
    ID: stream.ID,
    StreamName: stream.StreamName,
    Buffs: bs,
  }

  i.streams[id] = &ns

  return &ns, nil
}

func (i *inMemStore) ListStreams() ([]*StreamList, error) {
  sl := []*StreamList{}

  for _, s := range i.streams {
    sl = append(sl, &StreamList{
      ID: s.ID,
      StreamName: s.StreamName,
    })
  }

  return sl, nil
}

func NewInMemStore() Store {
	return &inMemStore{
		mu:        &sync.RWMutex{},
		buffCounter: 0,
		buffs:     make(map[uint64]*Buff),
		streamCounter: 0,
		streams:   make(map[uint64]*Stream),
	}
}
