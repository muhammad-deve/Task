-- name: LogDriverActivity :one
INSERT INTO driver_activity_log (
    driver_id,
    action,
    notes
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetDriverActivityLog :many
SELECT * FROM driver_activity_log
WHERE driver_id = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;

-- name: GetDriverWorkingHours :one
WITH activity_pairs AS (
    SELECT 
        driver_id,
        action,
        timestamp,
        LEAD(timestamp) OVER (PARTITION BY driver_id ORDER BY timestamp) as next_timestamp,
        LEAD(action) OVER (PARTITION BY driver_id ORDER BY timestamp) as next_action
    FROM driver_activity_log
    WHERE driver_id = $1
        AND timestamp >= $2
        AND timestamp <= $3
)
SELECT 
    COALESCE(SUM(
        CASE 
            WHEN action = 'went_online' AND next_action = 'went_offline'
            THEN EXTRACT(EPOCH FROM (next_timestamp - timestamp))
            ELSE 0
        END
    ), 0) as total_seconds
FROM activity_pairs;

-- name: GetActiveDriversCount :one
WITH latest_activity AS (
    SELECT DISTINCT ON (driver_id) 
        driver_id,
        action
    FROM driver_activity_log
    ORDER BY driver_id, timestamp DESC
)
SELECT COUNT(*) 
FROM latest_activity 
WHERE action = 'went_online';
