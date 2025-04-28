CREATE TABLE article_revisions(
  id bigserial not null primary key,
  article_id bigserial not null,
  user_id bigserial not null,
  user_ip varchar not null,
  title varchar not null unique,
  content text not null,
  is_deleted boolean default false,
  created_at timestamp,
  updated_at timestamp
);