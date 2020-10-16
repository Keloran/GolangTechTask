package handlers

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"

  "github.com/go-chi/chi"
)

type streamHandler struct {
  store StreamStore
}

func (s *streamHandler) ListStreams(w http.ResponseWriter, r *http.Request) {
  streams, err := s.store.ListStreams()
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to retrieve streams, %*v", err), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(streams); err != nil {
    http.Error(w, fmt.Sprintf("failed to encode streams, %+v", err), http.StatusInternalServerError)
    return
  }
  return
}

func (s *streamHandler) CreateStream(w http.ResponseWriter, r *http.Request) {
  cs := CreateStream{}
  if err := json.NewDecoder(r.Body).Decode(&cs); err != nil {
    http.Error(w, fmt.Sprintf("failed to unmarshall request, %+v", err), http.StatusBadRequest)
    return
  }

  stream := &Stream{
    ID: 0,
    StreamName: cs.StreamName,
  }
  id, err := s.store.SetStream(stream)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to create stream, %+v", err), http.StatusInternalServerError)
    return
  }

  stream, err = s.store.GetStream(id)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to get stream, %+v", err), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(stream); err != nil {
    http.Error(w, fmt.Sprintf("failed to encode stream, %+v", err), http.StatusInternalServerError)
    return
  }
}

func (s *streamHandler) AttachBuffs(w http.ResponseWriter, r *http.Request) {
  id := chi.URLParam(r, "id")
  if id == "" {
    http.Error(w, "missing stream id in path", http.StatusBadRequest)
    return
  }

  attach := StreamAttachBuffs{}
  if err := json.NewDecoder(r.Body).Decode(&attach); err != nil {
    http.Error(w, fmt.Sprintf("failed to unmarshall request, %+v", err), http.StatusBadRequest)
    return
  }

  parsedid, err := strconv.ParseUint(id, 10, 64)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to parse id, %v", err), http.StatusBadRequest)
    return
  }

  stream, err := s.store.AppendBuffs(parsedid, attach.Buffs)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to append buffs, %+v", err), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(stream); err != nil {
    http.Error(w, fmt.Sprintf("failed to encode strean, %+v", err), http.StatusInternalServerError)
    return
  }
}

func (s *streamHandler) GetStream(w http.ResponseWriter, r *http.Request) {
  id := chi.URLParam(r, "id")
  if id == "" {
    http.Error(w, "missing stream id in path", http.StatusBadRequest)
    return
  }

  parsedid, err := strconv.ParseUint(id, 10, 64)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to parse id, %v", err), http.StatusBadRequest)
    return
  }

  stream, err := s.store.GetStream(parsedid)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to get stream, %v", err), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(stream); err != nil {
    http.Error(w, fmt.Sprintf("failed encode stream, %v", err), http.StatusInternalServerError)
    return
  }
}
