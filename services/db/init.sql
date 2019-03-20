-- 
-- robots
-- 

CREATE DATABASE robotdb;
USE robotdb;

CREATE TABLE robots (
    serial text not null unique,
    nickname text,
    robotType text,
    registeredUserId text
);

INSERT INTO robots (serial, nickname, robotType, registeredUserId) VALUES ('123abc', 'testLightbot', 'switch', '1');
INSERT INTO robots (serial, nickname, robotType, registeredUserId) VALUES ('qwerty', 'testThermoBot', 'thermostat', '1');

-- 
-- switches
-- 

CREATE DATABASE switchdb;
USE switchdb;

CREATE TABLE switches (
    -- ID of the robot
    serial text not null unique,
    -- Weather the switch is on or off
    isOn boolean not null,
    -- Angles to set servo too
    onAngle int not null,
    offAngle int not null,
    restAngle int not null,
    -- Weather the robot has been calibrated
    isCalibrated boolean not null
);

INSERT INTO switches (serial, isOn, onAngle, offAngle, restAngle, isCalibrated) VALUES ('123abc', false, 90, 0, 45, true);

-- 
-- thermostats
-- 

CREATE DATABASE thermodb;
USE thermodb;

CREATE TABLE thermostats (
    -- ID of the robot
    serial text not null unique,
    -- Tempreture of current thermostat
    tempreture int not null,
    -- Angles to set servo too
    minAngle int not null,
    maxAngle int not null,
    -- Min/max of tempretures
    minTempreture int not null,
    maxTempreture int not null,
    -- Weather the robot has been calibrated
    isCalibrated boolean not null
);

INSERT INTO thermostats (serial, tempreture, minAngle, maxAngle, minTempreture, maxTempreture, isCalibrated) VALUES ('qwerty', 293, 30, 170, 283, 303, true);

-- 
-- users
-- 

CREATE DATABASE userdb;
USE userdb;

CREATE TABLE users (
    id serial,
    email text not null unique,
    password text not null 
);
