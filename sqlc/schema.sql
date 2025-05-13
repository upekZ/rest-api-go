CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_status AS ENUM ('Active', 'Inactive');

CREATE TABLE IF NOT EXISTS "user" (
		userId uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
		first_name varchar(100) NOT NULL,
		last_name varchar(100) NOT NULL,
		email varchar(100) NOT NULL,
		phone varchar(100),
		age integer,
		"status" user_status DEFAULT 'Active'
	);