-- # Jun 25 2020
create table if not exists users
(
    userId    int auto_increment
        primary key,
    fullName  varchar(90)  null,
    email     varchar(45)  null,
    password  varchar(255) null,
    createdAt int          null,
    updatedAt int          null
);

