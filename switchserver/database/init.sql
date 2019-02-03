CREATE TABLE switches (
    -- ID of the robot
    robotId int not null unique,
    -- Weather the switch is on or off
    isOn boolean not null,
    -- Angles to set servo too
    onAngle int not null,
    offAngle int not null,
    restAngle int not null,
    -- Weather the robot has been calibrated
    isCalibrated boolean not null
);

-- TODO: move into test database
INSERT INTO switches (robotId, isOn, onAngle, offAngle, restAngle, isCalibrated) VALUES (9999, false, 90, 0, 45, true);
