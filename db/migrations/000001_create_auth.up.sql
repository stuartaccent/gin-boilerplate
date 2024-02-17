begin;

create table auth_users
(
    id              uuid                     default gen_random_uuid() not null primary key,
    email           varchar(320)                                       not null unique,
    hashed_password bytea                                              not null,
    first_name      varchar(120)                                       not null,
    last_name       varchar(120)                                       not null,
    is_active       boolean                  default true              not null,
    is_verified     boolean                  default false             not null,
    created_at      timestamp with time zone default clock_timestamp(),
    updated_at      timestamp with time zone default clock_timestamp()
);

create or replace function set_updated_at()
    returns trigger as
$$
begin
    NEW.updated_at = clock_timestamp();
    return NEW;
end;
$$ language plpgsql;

create trigger set_updated_at
    before update
    on auth_users
    for each row
execute procedure set_updated_at();

commit;