drop table if exists users;
drop table if exists logins;

create table users (
    id serial primary key,
    firstName varchar(200),
    lastName varchar(200)
);

create table logins (
    id serial primary key,
    username varchar(100) unique not null,
    email varchar(254) unique not null,
    passwordHash text not null,
    created timestamp not null,
    updated timestamp not null,
    constraint userid
        foreign key(id)
            references users(id)
                on delete cascade
);


