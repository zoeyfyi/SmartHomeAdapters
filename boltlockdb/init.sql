CREATE TABLE boltlocks (
    -- ID of the robot
    serial text not null unique,
    -- Weather the boltlock is on or off
    isOn boolean not null,
    -- Angles to set servo too
    onAngle int not null,
    offAngle int not null,
    -- Weather the robot has been calibrated
    isCalibrated boolean not null
);

INSERT INTO boltlocks (serial, isOn, onAngle, offAngle, isCalibrated) VALUES ('bolty', false, 90, 0, true);
