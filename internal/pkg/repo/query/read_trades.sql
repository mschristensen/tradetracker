SELECT id, instrument_id, price, size, timestamp
FROM trades
WHERE instrument_id=$1::bigint AND timestamp > to_timestamp($2::bigint) AT TIME ZONE 'UTC'
ORDER BY timestamp ASC;
