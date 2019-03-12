CREATE TABLE command (
    -- ID of the robot
    serial text not null unique,
    -- Angle to turn the servo
    angle int not null,
    -- Weather it has been completed or not
    isCompleted boolean not null,
    -- time command was submitted
    submitTime datetime not null,
    -- time command was completed
    completeTime datetime not null,
    -- time between the commands (microseconds)
    delayTime int not null,
);