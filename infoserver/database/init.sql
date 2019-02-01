CREATE TABLE toggleRobots (
    serial text not null unique,
    nickname text,
    robotType text,
    stateOn text
);
CREATE TABLE rangeRobots (
    serial text not null unique,
    nickname text,
    robotType text,
    minimum integer,
    maximum integer
);
INSERT INTO toggleRobots (serial, nickname, robotType, stateOn) VALUES ('123abc', 'testLightbot', 'switch', false);
INSERT INTO rangeRobots (serial, nickname, robotType, minimum, maximum) VALUES ('T2D2', 'testThermoBot', 'thermostat', 0, 100);
