CREATE TABLE robots (
    id serial not null unique,
    nickname text,
    robotType text not null,
    minimum integer,
    maximum integer
);
