alter table posts 
add column likes_count int not null default 0,
add column views_count int not null default 0,
add column replies_count int not null default 0;
