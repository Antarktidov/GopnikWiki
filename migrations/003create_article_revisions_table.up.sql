CREATE TABLE article_revisions(
  id bigserial not null primary key,
  article_id bigserial not null,
  user_id bigserial not null,
  user_ip varchar not null,
  title varchar not null unique,
  content text not null unique,
  is_deleted boolean default false
);