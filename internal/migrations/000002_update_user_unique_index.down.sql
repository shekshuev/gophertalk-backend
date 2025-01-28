drop index idx__users__user_name;
alter table users add constraint uk__users__user_name unique(user_name)