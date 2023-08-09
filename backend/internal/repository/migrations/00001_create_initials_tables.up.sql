CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_uid text UNIQUE,
    simulated BOOLEAN,
    name text
);

CREATE INDEX IF NOT EXISTS users_user_uid_index ON users(user_uid);

CREATE TABLE IF NOT EXISTS vehicles (
    id SERIAL PRIMARY KEY,
    registration_country text,
    registration_number text,
    owner_id int references users(id),
    icon text,
    UNIQUE(registration_country, registration_number)
);

CREATE TABLE IF NOT EXISTS vehicle_positions (
    id SERIAL PRIMARY KEY,
    vehicle_id int references vehicles(id),
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    recorded_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS ride_requests (
    id SERIAL PRIMARY KEY,
    rider_id int references users(id),
    driver_id int null references users(id),
    from_lat DOUBLE PRECISION,
    from_lng DOUBLE PRECISION,
    from_name text,
    to_lat DOUBLE PRECISION,
    to_lng DOUBLE PRECISION,
    to_name text,
    state int,
    directions_json_version int null,
    directions_json text null,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);