create table if not exists users (
    id bigserial,
    user_name varchar(30) not null,
    first_name varchar(30) not null,
    last_name varchar(30) not null,
    password_hash varchar(72) not null,
    status smallint default 1,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    constraint pk__users primary key(id),
    constraint uk__users__user_name unique(user_name),
    constraint chk__users__status check(status in (0, 1))
)