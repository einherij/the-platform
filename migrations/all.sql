
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL
);

CREATE TABLE balances (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    balance INT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO users (username) VALUES ('test_name_1');
INSERT INTO users (username) VALUES ('test_name_2');
INSERT INTO balances (user_id, balance) VALUES ((SELECT id FROM users WHERE username='test_name_1'), 100);
INSERT INTO balances (user_id, balance) VALUES ((SELECT id FROM users WHERE username='test_name_2'), 200);
