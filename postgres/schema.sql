drop table if exists users;
drop table if exists logins;

create table logins (
    id serial primary key,
    username varchar(100) unique not null,
    email varchar(254) unique not null,
    passwordHash text not null,
    passwordSalt text not null,
    created timestamp not null,
    updated timestamp not null
);

create table users (
    id serial primary key,
    loginID int references logins(id) not null,
);

