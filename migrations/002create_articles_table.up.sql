CREATE TABLE articles (
  id bigserial not null primary key,
  title varchar not null unique,
  is_deleted boolean default false,
  created_at timestamp,
  updated_at timestamp
);