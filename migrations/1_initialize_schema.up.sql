begin;
create table organisations(
    id serial primary key,
    name text not null unique,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    deleted_at timestamptz
);
create table users(
    id serial primary key,
    username text not null unique,
    avatar text default '',
    password_hash text not null,
    api_key text not null,
    created_at timestamptz not null default current_timestamp,
    deleted_at timestamptz
);
create table organisation_comments(
    id serial primary key,
    comment text not null,
    organisation_id integer not null references organisations (id),
    created_by integer not null references users (id),
    created_at timestamptz not null default current_timestamp,
    deleted_at timestamptz
);
create table user_organisations(
    user_id integer references users (id),
    organisation_id integer references organisations (id),
    created_at timestamptz not null default current_timestamp,
    unique(user_id, organisation_id)
);
create table user_followers(
    followee_id integer references users (id),
    follower_id integer references users (id),
    created_at timestamptz not null default current_timestamp,
    check(followee_id != follower_id),
    unique(follower_id, followee_id)
);
commit;
