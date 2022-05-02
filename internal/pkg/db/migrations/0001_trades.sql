-- +migrate Up
CREATE TABLE trades (
  id serial PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  instrument_id bigint NOT NULL,
  size bigint NOT NULL,
  price numeric NOT NULL,
  timestamp timestamp without time zone NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS trades;
