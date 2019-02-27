CREATE TABLE users (
    id serial,
    email text not null unique,
    password text not null 
);
