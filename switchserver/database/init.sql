CREATE TABLE switches (
    -- ID of the robot
    robotId int not null unique,
    -- Weather the switch is on or off
    isOn boolean not null
);
