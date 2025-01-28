create table if not exists likes (
    user_id bigint,
    post_id bigint,
    created_at timestamp not null default now(),
    constraint pk__likes primary key (user_id, post_id),
    constraint fk__likes__user_id foreign key (user_id) references users(id),
    constraint fk__likes__post_id foreign key (post_id) references posts(id)
);

create table if not exists views (
    user_id bigint,
    post_id bigint,
    created_at timestamp not null default now(),
    constraint pk__views primary key (user_id, post_id),
    constraint fk__views__user_id foreign key (user_id) references users(id),
    constraint fk__views__post_id foreign key (post_id) references posts(id)
);