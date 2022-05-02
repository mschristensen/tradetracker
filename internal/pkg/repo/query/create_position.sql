INSERT INTO positions (instrument_id, size, timestamp)
VALUES ($1::int, $2::int, to_timestamp($3::bigint) AT TIME ZONE 'UTC')
RETURNING id;
