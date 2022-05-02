DELETE FROM positions
WHERE instrument_id=$1::bigint;
