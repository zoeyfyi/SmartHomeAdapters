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
