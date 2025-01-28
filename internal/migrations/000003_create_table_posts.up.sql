create table if not exists posts (
    id bigserial,
    text varchar(280) not null,
    repost_of_id bigint,
    user_id bigint,
    created_at timestamp not null default now(),
    deleted_at timestamp,
    constraint pk__posts primary key(id),
    constraint fk__posts__user_id foreign key(user_id) references users(id),
    constraint fk__posts__repost_of_id foreign key(repost_of_id) references posts(id)
)