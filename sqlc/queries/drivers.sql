-- name: CreateDriver :one
INSERT INTO drivers (
    full_name,
    phone,
    license_number,
    car_model,
    car_plate,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetDriverByID :one
SELECT * FROM drivers 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetDrivers :many
SELECT * FROM drivers
WHERE deleted_at IS NULL
AND ($1::text IS NULL OR status = $1)
AND ($2::text IS NULL OR full_name ILIKE '%' || $2 || '%' OR phone ILIKE '%' || $2 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountDrivers :one
SELECT COUNT(*) FROM drivers
WHERE deleted_at IS NULL
AND ($1::text IS NULL OR status = $1)
AND ($2::text IS NULL OR full_name ILIKE '%' || $2 || '%' OR phone ILIKE '%' || $2 || '%');

-- name: UpdateDriver :one
UPDATE drivers
SET
    full_name = COALESCE($2, full_name),
    phone = COALESCE($3, phone),
    license_number = COALESCE($4, license_number),
    car_model = COALESCE($5, car_model),
    car_plate = COALESCE($6, car_plate),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateDriverStatus :one
UPDATE drivers
SET
    status = $2,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteDriver :exec
UPDATE drivers
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CheckPhoneExists :one
SELECT EXISTS(
    SELECT 1 FROM drivers 
    WHERE phone = $1 AND deleted_at IS NULL AND id != COALESCE($2::uuid, '00000000-0000-0000-0000-000000000000'::uuid)
);

-- name: CheckLicenseExists :one
SELECT EXISTS(
    SELECT 1 FROM drivers 
    WHERE license_number = $1 AND deleted_at IS NULL AND id != COALESCE($2::uuid, '00000000-0000-0000-0000-000000000000'::uuid)
);

-- name: CheckCarPlateExists :one
SELECT EXISTS(
    SELECT 1 FROM drivers 
    WHERE car_plate = $1 AND deleted_at IS NULL AND id != COALESCE($2::uuid, '00000000-0000-0000-0000-000000000000'::uuid)
);
