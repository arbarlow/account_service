CREATE TABLE account_service.accounts_map_id (
    id text PRIMARY KEY,
    createdat timestamp,
    email text,
    hashedpassword text,
    name text
);

CREATE TABLE account_service.accounts_map_email (
    email text PRIMARY KEY,
    id text
);
