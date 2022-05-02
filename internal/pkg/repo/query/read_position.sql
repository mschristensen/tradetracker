SELECT id, instrument_id, size, timestamp
FROM positions
WHERE instrument_id=$1::bigint
AND timestamp <= to_timestamp($2::bigint) AT TIME ZONE 'UTC'
ORDER BY timestamp DESC
LIMIT 1::bigint;
