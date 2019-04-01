CREATE TABLE boolparameter (
    -- ID of parameter
    serial text not null unique,
    -- robot ID
    robotId text not null,    
    value boolean not null
);

CREATE TABLE intparameter (
    -- ID of parameter
    serial text not null unique,
    -- robot ID
    robotId text not null,    
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