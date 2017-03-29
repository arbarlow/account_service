ALTER TABLE accounts ADD COLUMN confirmation_token text;
ALTER TABLE accounts ADD COLUMN password_reset_token text;


CREATE UNIQUE INDEX IF NOT EXISTS accounts_confirmation_token ON accounts (confirmation_token);
CREATE UNIQUE INDEX IF NOT EXISTS accounts_reset_token ON accounts (password_reset_token);
