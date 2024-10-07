create table if not exists users (
    id serial primary key,
    firstname varchar(50) not null,
    lastname varchar(50) not null,
    username varchar(50) not null unique,
    password varchar(100) not null
);