CREATE TABLE driver_activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id UUID NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    action TEXT NOT NULL CHECK (action IN ('went_online', 'went_offline')),
    timestamp TIMESTAMP NOT NULL DEFAULT now(),
    notes TEXT
);

CREATE INDEX idx_driver_activity_driver_id ON driver_activity_log(driver_id);
CREATE INDEX idx_driver_activity_timestamp ON driver_activity_log(timestamp DESC);
