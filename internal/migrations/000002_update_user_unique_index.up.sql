alter table users drop constraint uk__users__user_name;
create unique index idx__users__user_name on users(user_name) where (deleted_at is null);