CREATE TABLE drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL CHECK (length(full_name) BETWEEN 3 AND 100),
    phone TEXT NOT NULL UNIQUE CHECK (phone ~ '^\+998[0-9]{9}$'),
    license_number TEXT NOT NULL UNIQUE,
    car_model TEXT NOT NULL,
    car_plate TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'blocked')),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_drivers_phone ON drivers(phone);
CREATE INDEX idx_drivers_status ON drivers(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_drivers_deleted_at ON drivers(deleted_at);
