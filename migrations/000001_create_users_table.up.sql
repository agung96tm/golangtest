CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name varchar(200) NOT NULL,
    email varchar(200) NOT NULL,
    password varchar(100) NOT NULL,
    role_id int,
    role_name varchar(200),
    last_access timestamp(0) with time zone NOT NULL DEFAULT NOW()
);