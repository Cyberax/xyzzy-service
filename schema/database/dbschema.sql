-- Enable geo features
create extension if not exists "postgis";
create extension if not exists "uuid-ossp";

create table test(id uuid primary key default uuid_generate_v4(),
    name varchar not null, url varchar);

create schema if not exists sqlc;
CREATE or replace FUNCTION sqlc.arg(val varchar) RETURNS record AS $$
BEGIN
   RETURN val + 1;
END; $$ language plpgsql;
