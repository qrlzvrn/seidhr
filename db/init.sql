CREATE EXTENSION pg_trgm;

CREATE TABLE medicament (
    id SERIAL PRIMARY KEY,
    title VARCHAR(128),
    availability BOOLEAN
);

CREATE TABLE tguser (
    id INTEGER PRIMARY KEY,
    chat_id INTEGER,
    state VARCHAR(30),
    selected_med INTEGER
);

CREATE TABLE subscription (
    medicament_id INTEGER REFERENCES medicament(id),
    tguser_id INTEGER REFERENCES tguser(id)
);