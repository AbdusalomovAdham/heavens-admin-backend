-- users --
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    full_name text NOT NULL,
    password text NOT NULL,
    email text NOT NULL UNIQUE,
    role integer DEFAULT 3 NOT NULL,
    district_id integer NOT NULL,
    region_id integer NOT NULL,
    status boolean DEFAULT false NOT NULL,
    avatar text,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp,
    created_by int,
    updated_by int,
    );

ALTER  TABLE  users
    owner TO postgres;

-- regions --
CREATE TABLE IF NOT EXISTS regions(
    id serial PRIMARY KEY,
    region text NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    created_by int,
    updated_by int DEFAULT NULL,
    updated_at timestamp,
    );

ALTER  TABLE  regions
    owner TO postgres;


-- distracts --
CREATE TABLE IF NOT EXISTS distracts(
    id serial PRIMARY KEY,
    distract text NOT NULL,
    region_id int NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    created_by int,
    updated_by int DEFAULT NULL,
    updated_at timestamp,
    );

ALTER  TABLE  distracts
    owner TO postgres;
