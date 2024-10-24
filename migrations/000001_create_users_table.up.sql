create table if not exists users (
    id serial primary key,
    firstname varchar(255) not null,
    lastname varchar(255) not null,
    username varchar(255) not null unique,
    password varchar(255) not null
);