CREATE TABLE boolparameter (
    -- ID of parameter
    serial text not null,
    -- robot ID
    robotId text not null,
    PRIMARY KEY(serial, robotId),  
    value boolean not null
);

CREATE TABLE intparameter (
    -- ID of parameter
    serial text not null,
    -- robot ID
    robotId text not null,
    PRIMARY KEY(serial, robotId),
    value int not null
);

CREATE TABLE togglestatus (
    robotId text not null unique,
    value boolean not null
);

CREATE TABLE rangestatus (
    robotId text not null unique,
    value int not null
);