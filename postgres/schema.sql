drop table if exists users;
drop table if exists logins;

create table logins (
    id serial primary key,
    username varchar(100) unique not null,
    email varchar(254) unique not null,
    passwordHash text not null,
    created timestamp not null,
    updated timestamp not null
);

create table users (
    id serial primary key,
    firstName varchar(200),
    lastName varchar(200), 
    phone varchar(15),
    created timestamp not null,
    updated timestamp not null,
    loginid integer unique not null,
    constraint fk_loginid
        foreign key(loginid)
            references logins(id)
                on delete cascade
);


