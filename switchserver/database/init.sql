CREATE TABLE switches (
    -- ID of the robot
    robotId serial,
    -- True when the switch is on, false otherwise
    switchOn boolean not null
);
