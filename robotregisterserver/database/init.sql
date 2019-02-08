CREATE TABLE toggleRobots (
    serial text not null unique,
    nickname text,
    robotType text,
    registeredUserId text
);
CREATE TABLE rangeRobots (
    serial text not null unique,
    nickname text,
    robotType text,
    minimum integer,
    maximum integer,
    registeredUserId text
);
INSERT INTO toggleRobots (serial, nickname, robotType, registeredUserId) VALUES ('123abc', 'testLightbot', 'switch', '1');
INSERT INTO rangeRobots (serial, nickname, robotType, minimum, maximum, registeredUserId) VALUES ('T2D2', 'testThermoBot', 'thermostat', 0, 100, '1');
