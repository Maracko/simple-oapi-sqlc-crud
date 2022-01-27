CREATE TABLE todos (
  id   BIGSERIAL PRIMARY KEY,
  title TEXT,
  tags  VARCHAR(255)[],
  content TEXT
);