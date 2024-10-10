CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    text TEXT,
    release_date VARCHAR(50),
    link VARCHAR(255)
);