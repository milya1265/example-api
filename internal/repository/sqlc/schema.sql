CREATE TABLE users
(
    id  varchar(40) unique PRIMARY KEY,
    login varchar(30) unique not null ,
    password_hash varchar(30)  not null
);

CREATE TABLE jwt
(
    user_id           varchar(40) unique  not null ,
    foreign key (user_id) REFERENCES users (id),
    access_token varchar(200)  not null,
    refresh_token varchar(200) not null
);

CREATE TABLE users_role
(
    user_id   varchar(40) unique  not null ,
    foreign key (user_id) REFERENCES users (id),
    role integer  not null
);


