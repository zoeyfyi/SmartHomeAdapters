CREATE TABLE robots (
    serial text not null unique,
    nickname text,
    robotType text,
    registeredUserId text
);

INSERT INTO robots (serial, nickname, robotType, registeredUserId) VALUES ('123abc', 'testLightbot', 'switch', '1');
INSERT INTO robots (serial, nickname, robotType, registeredUserId) VALUES ('qwerty', 'testThermoBot', 'thermostat', '1');