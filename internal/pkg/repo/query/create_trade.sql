INSERT INTO trades (instrument_id, size, price, timestamp)
VALUES ($1::int, $2::int, $3::numeric, to_timestamp($4::bigint) AT TIME ZONE 'UTC')
RETURNING id;
