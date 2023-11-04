CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SET TIME ZONE 'Europe/Istanbul';

CREATE TABLE users (
		id INT NOT NULL SERIAL,
    uuid UUID DEFAULT uuid_generate_v4() ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    user_name VARCHAR(25) NOT NULL,
    user_surname VARCHAR(25) NOT NULL,
    nickname varchar(25) NOT NULL,
    email VARCHAR (100) NOT NULL,
    password_hash VARCHAR (255) NOT NULL,
    user_status INT NOT NULL,
    user_role VARCHAR (25) NOT NULL
    PRIMARY KEY (id)
    UNIQUE (nickname)
    UNIQUE (email)
    UNIQUE  (uuid)
);


CREATE TABLE user_friends (
    id int NOT NULL AUTO_INCREMENT,
    created_at datetime(3) DEFAULT NULL,
    updated_at datetime(3) DEFAULT NULL,
    deleted_at bigint unsigned DEFAULT NULL,
    user_id int DEFAULT NULL,
    friend_id int DEFAULT NULL,
    PRIMARY KEY (id),
);