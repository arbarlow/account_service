CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS accounts (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	name text NULL,
	email text NOT NULL,
	hashed_password text NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE UNIQUE INDEX IF NOT EXISTS accounts_email ON accounts ((lower(email)));
