CREATE TABLE users (
  id bigserial not null primary key,
  username varchar not null unique,
  is_admin boolean default false,
  encrypted_password varchar not null,
  created_at timestamp
);