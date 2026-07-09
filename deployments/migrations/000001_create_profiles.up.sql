-- deployments/migrations/000001_create_profiles.up.sql

CREATE TABLE profiles
(
    id CHAR(26) PRIMARY KEY,

    full_name TEXT NOT NULL,

    headline TEXT,

    bio TEXT,

    email TEXT,

    phone TEXT,

    location TEXT,

    github TEXT,

    linkedin TEXT,

    website TEXT,

    avatar TEXT,

    resume_url TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);