package handlers

import (
  "context"
  "fmt"
  "os"

  "github.com/jackc/pgx/v4"
)

type storage struct {
}

func (i *storage) GetBuff(id uint64) (*Buff, error) {
  db, err := pgx.Connect(context.Background(), fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE")))
  if err != nil {
    fmt.Printf("getBuff connection failed, %+v", err)
  }
  defer func() {
    err := db.Close(context.Background())
    if err != nil {
      fmt.Printf("getBuff dbClose: %+v", err)
    }
  }()

  data, err := db.Query(context.Background(), "SELECT question, option FROM buff LEFT JOIN buff_options ON buff_options.buff_id = buff.id WHERE buff.id = $1", id)
  if err != nil {
    return nil, fmt.Errorf("buff not found, %w", err)
  }

  b := &Buff{
    ID: id,
  }
  var option, question string
  for data.Next() {
    err := data.Scan(&question, &option)
    if err != nil {
      return b, fmt.Errorf("row failed to scan, %w", err)
    }
    b.Answers = append(b.Answers, option)
  }
  b.Question = question

  return b, nil
}

func (i *storage) SetBuff(b *Buff) (uint64, error) {
  db, err := pgx.Connect(context.Background(), fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE")))
  if err != nil {
    fmt.Printf("setBuff connection failed, %+v", err)
  }
  defer func() {
    err := db.Close(context.Background())
    if err != nil {
      fmt.Printf("setBuff dbClose: %+v", err)
    }
  }()

  var insertId uint64
	err = db.QueryRow(context.Background(), "INSERT INTO buff (question) VALUES ($1) RETURNING (id)", b.Question).Scan(&insertId)
	if err != nil {
	  return 0, fmt.Errorf("failed to insert buff, %w", err)
  }

  for _, a := range b.Answers {
    _, err = db.Exec(context.Background(), "INSERT INTO buff_options (buff_id, option) VALUES ($1, $2)", insertId, a)
    if err != nil {
      return insertId, fmt.Errorf("failed to insert option, %w", err)
    }
  }

	return insertId, nil
}

func (i *storage) ListBuffs() (map[uint64]*Buff, error) {
  db, err := pgx.Connect(context.Background(), fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE")))
  if err != nil {
    fmt.Printf("listBuff sql connection failed, %+v", err)
  }
  defer func() {
    err := db.Close(context.Background())
    if err != nil {
      fmt.Printf("listBuff dbClose: %+v", err)
    }
  }()

  data, err := db.Query(context.Background(), "SELECT id, question FROM buff")
  if err != nil {
    return nil, fmt.Errorf("list buffs, %w", err)
  }

  var incId uint64 = 0
  buffs := make(map[uint64]*Buff)
  for data.Next() {
    b := Buff{}

    err := data.Scan(&b.ID, &b.Question)
    if err != nil {
      return buffs, fmt.Errorf("list buffs scan, %w", err)
    }
    buffs[incId] = &b
    incId++
  }

  return buffs, nil
}

func (i *storage) GetStream(id uint64) (*Stream, error) {
  db, err := pgx.Connect(context.Background(), fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE")))
  if err != nil {
    fmt.Printf("getStream connection failed, %+v", err)
  }
  defer func() {
    err := db.Close(context.Background())
    if err != nil {
      fmt.Printf("getStream dbClose: %+v", err)
    }
  }()

  s := &Stream{
    ID: id,
  }

  err = db.QueryRow(context.Background(), "SELECT name AS stream_name FROM stream WHERE stream.id = $1", id).Scan(&s.StreamName)
  if err != nil {
    return nil, fmt.Errorf("stream not found, %w", err)
  }

  data, err := db.Query(context.Background(), "SELECT buff_id FROM stream_buffs WHERE stream_id = $1", id)
  if err != nil {
    return s, fmt.Errorf("get stream buffs, %w", err)
  }
  for data.Next() {
    var buffId uint64 = 0
    if err := data.Scan(&buffId); err != nil {
      return s, fmt.Errorf("buffid scan, %w", err)
    }
    b, err := i.GetBuff(buffId)
    if err != nil {
      return s, fmt.Errorf("getStream buff, %w", err)
    }
    s.Buffs = append(s.Buffs, *b)
  }

  return s, nil
}

func (i *storage) SetStream(s *Stream) (uint64, error) {
  db, err := pgx.Connect(context.Background(), fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE")))
  if err != nil {
    fmt.Printf("setStream connection failed, %+v", err)
  }
  defer func() {
    err := db.Close(context.Background())
    if err != nil {
      fmt.Printf("setStream dbClose: %+v", err)
    }
  }()

  var insertId uint64
  err = db.QueryRow(context.Background(), "INSERT INTO stream (name) VALUES ($1) RETURNING (id)", s.StreamName).Scan(&insertId)
  if err != nil {
    return 0, fmt.Errorf("failed to insert stream, %w", err)
  }
  return insertId, nil
}

func (i *storage) AppendBuffs(id uint64, buffs []uint64) (*Stream, error) {
  db, err := pgx.Connect(context.Background(), fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE")))
  if err != nil {
    fmt.Printf("appendBuffs connection failed, %+v", err)
  }
  defer func() {
    err := db.Close(context.Background())
    if err != nil {
      fmt.Printf("appendBuffs dbClose: %+v", err)
    }
  }()

  for _, b := range buffs {
    _, err := db.Exec(context.Background(), "INSERT INTO stream_buffs (stream_id, buff_id) VALUES ($1, $2)", id, b)
    if err != nil {
      return nil, fmt.Errorf("failed to insert stream buffs, %w", err)
    }
  }
  return i.GetStream(id)
}

func (i *storage) ListStreams() ([]*StreamList, error) {
  db, err := pgx.Connect(context.Background(), fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USERNAME"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_DATABASE")))
  if err != nil {
    fmt.Printf("listStreams connection failed, %+v", err)
  }
  defer func() {
    err := db.Close(context.Background())
    if err != nil {
      fmt.Printf("listStreams dbClose: %+v", err)
    }
  }()

  sl := []*StreamList{}

  data, err := db.Query(context.Background(), "SELECT id, name FROM stream")
  if err != nil {
    return nil, fmt.Errorf("stream list not found, %w", err)
  }
  for data.Next() {
    s := StreamList{}
    err := data.Scan(&s.ID, &s.StreamName)
    if err != nil {
      return sl, fmt.Errorf("failed to scan streamlist, %w", err)
    }

    sl = append(sl, &s)
  }

  return sl, nil
}

func Storage() Store {
  return &storage{}
}
