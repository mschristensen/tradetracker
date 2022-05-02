-- +migrate Up
CREATE TABLE trades (
  id serial PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  instrument_id int NOT NULL,
  size int NOT NULL,
  price money NOT NULL,
  timestamp timestamp without time zone NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS trades;
