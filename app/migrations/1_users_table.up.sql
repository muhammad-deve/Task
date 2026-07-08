CREATE TABLE "users" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "fullName" TEXT,
    "phoneNumber" TEXT UNIQUE,
    "passwordHash" TEXT,
    "createdAt" TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
    "updatedAt" TIMESTAMP WITHOUT TIME ZONE,
    "deletedAt" TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS "idx_users_phone" ON "users" ("phoneNumber");