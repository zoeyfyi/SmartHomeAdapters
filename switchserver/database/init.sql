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
