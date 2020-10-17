CREATE TABLE buff (
    id bigserial,
    question text,
    PRIMARY KEY (id)
);

CREATE TABLE buff_options (
    id bigserial,
    buff_id int NOT NULL,
    option text,
    PRIMARY KEY (id),
    FOREIGN KEY (buff_id) REFERENCES buff(id) ON DELETE CASCADE
);

CREATE TABLE stream (
    id bigserial,
    name text,
    PRIMARY KEY (id)
);

CREATE TABLE stream_buffs (
    stream_id int NOT NULL,
    buff_id int NOT NULL,
    PRIMARY KEY (stream_id, buff_id),
    FOREIGN KEY (stream_id) REFERENCES stream(id) ON UPDATE CASCADE,
    FOREIGN KEY (buff_id) REFERENCES buff(id) ON UPDATE CASCADE
);
