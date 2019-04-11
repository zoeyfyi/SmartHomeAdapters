CREATE TABLE robots (
    serial text not null unique,
    nickname text,
    robotType text,
    registeredUserId text
);

INSERT INTO robots (serial) VALUES ('123abc');
INSERT INTO robots (serial) VALUES ('qwerty');