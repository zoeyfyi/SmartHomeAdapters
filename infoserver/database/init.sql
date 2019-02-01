CREATE TABLE robots (
    serial text not null unique,
    nickname text,
    robotType text not null,
    minimum integer,
    maximum integer
);
INSERT INTO robots (serial, nickname, robotType, minimum, maximum) VALUES ('123abc', 'testRobot', 'testType', 0, 100);
