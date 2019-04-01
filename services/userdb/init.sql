CREATE TABLE users (
    id serial,
    username text not null
    email text not null unique,
    password text not null 
);
