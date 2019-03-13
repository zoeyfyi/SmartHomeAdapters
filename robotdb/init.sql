CREATE TABLE command (
    -- ID of the robot
    id serial,
    -- Angle to turn the servo
    angle int not null,
    -- Weather it has been completed or not
    isCompleted boolean not null,
    -- time command was submitted
    submitTime timestamp not null,
    -- time command was completed
    completeTime timestamp,
    -- time between the commands (microseconds)
    delayTime int not null
);