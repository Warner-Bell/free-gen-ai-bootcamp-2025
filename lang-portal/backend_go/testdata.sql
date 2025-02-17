INSERT INTO groups (name, description, created_at, updated_at) 
VALUES 
('Basic Phrases', 'Common everyday phrases', DATETIME('now'), DATETIME('now')),
('Greetings', 'Basic greeting phrases', DATETIME('now'), DATETIME('now'));

INSERT INTO words (word, translation, notes, japanese, romaji, english, created_at, updated_at) 
VALUES 
('こんにちは', 'Hello', 'Formal greeting', 'こんにちは', 'konnichiwa', 'hello', DATETIME('now'), DATETIME('now')),
('ありがとう', 'Thank you', 'Basic thanks', 'ありがとう', 'arigatou', 'thank you', DATETIME('now'), DATETIME('now'));
