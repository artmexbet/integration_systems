CREATE TABLE IF NOT EXISTS laureates (
    id INT PRIMARY KEY,
    firstname VARCHAR(100) NOT NULL,
    surname VARCHAR(100),
    motivation TEXT NOT NULL,
    share INT NOT NULL
);

CREATE TABLE IF NOT EXISTS prizes (
    id SERIAL PRIMARY KEY,
    year INT NOT NULL,
    category VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS prizes_to_laureates (
    prize_id INT REFERENCES prizes(id) ON DELETE CASCADE,
    laureate_id INT REFERENCES laureates(id) ON DELETE CASCADE,
    PRIMARY KEY (prize_id, laureate_id)
);
