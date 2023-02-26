CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR,
    "last_name" VARCHAR,
    "email" VARCHAR NOT NULL,
    "password" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "urls" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "original_url" TEXT NOT NULL,
    "hashed_url" VARCHAR NOT NULL,
    "max_clicks" INT CHECK ("max_clicks" > 0),
    "expires_at" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);