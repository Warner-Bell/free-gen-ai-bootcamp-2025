-- Clear existing data
DELETE FROM word_review_items;
DELETE FROM study_sessions;
DELETE FROM words_groups;
DELETE FROM words;
DELETE FROM groups;
DELETE FROM study_activities;

-- Insert groups
INSERT INTO groups (name, description, created_at, updated_at) 
VALUES 
('Basic Phrases', 'Common everyday phrases', DATETIME('now'), DATETIME('now')),
('Greetings', 'Basic greeting phrases', DATETIME('now'), DATETIME('now')),
('Numbers', 'Basic numbers 1-10', DATETIME('now'), DATETIME('now')),
('Colors', 'Basic colors', DATETIME('now'), DATETIME('now'));

-- Insert words
INSERT INTO words (word, translation, notes, japanese, romaji, english, created_at, updated_at) 
VALUES 
('こんにちは', 'Hello', 'Formal greeting', 'こんにちは', 'konnichiwa', 'hello', DATETIME('now'), DATETIME('now')),
('ありがとう', 'Thank you', 'Basic thanks', 'ありがとう', 'arigatou', 'thank you', DATETIME('now'), DATETIME('now')),
('さようなら', 'Goodbye', 'Formal farewell', 'さようなら', 'sayounara', 'goodbye', DATETIME('now'), DATETIME('now')),
('おはよう', 'Good morning', 'Casual morning greeting', 'おはよう', 'ohayou', 'good morning', DATETIME('now'), DATETIME('now')),
('一', 'One', 'Number 1', '一', 'ichi', 'one', DATETIME('now'), DATETIME('now')),
('二', 'Two', 'Number 2', '二', 'ni', 'two', DATETIME('now'), DATETIME('now')),
('赤', 'Red', 'Color red', '赤', 'aka', 'red', DATETIME('now'), DATETIME('now')),
('青', 'Blue', 'Color blue', '青', 'ao', 'blue', DATETIME('now'), DATETIME('now'));

-- Insert words_groups relationships
INSERT INTO words_groups (word_id, group_id) 
VALUES 
(1, 1), (1, 2),  -- こんにちは belongs to both Basic Phrases and Greetings
(2, 1),          -- ありがとう in Basic Phrases
(3, 1), (3, 2),  -- さようなら in both Basic Phrases and Greetings
(4, 2),          -- おはよう in Greetings
(5, 3), (6, 3),  -- Numbers
(7, 4), (8, 4);  -- Colors

-- Insert study activities
INSERT INTO study_activities (name, description, created_at, updated_at)
VALUES 
('Vocabulary Review', 'Review and memorize vocabulary', DATETIME('now'), DATETIME('now')),
('Pronunciation Practice', 'Practice correct pronunciation', DATETIME('now'), DATETIME('now'));

-- Insert sample study sessions
INSERT INTO study_sessions (activity_id, started_at, ended_at, created_at, updated_at)
VALUES 
(1, DATETIME('now', '-1 hour'), DATETIME('now'), DATETIME('now'), DATETIME('now')),
(2, DATETIME('now', '-2 hour'), DATETIME('now', '-1 hour'), DATETIME('now'), DATETIME('now'));

-- Insert sample word reviews
INSERT INTO word_review_items (word_id, session_id, score, reviewed_at, created_at, updated_at)
VALUES 
(1, 1, 5, DATETIME('now'), DATETIME('now'), DATETIME('now')),
(2, 1, 4, DATETIME('now'), DATETIME('now'), DATETIME('now')),
(3, 2, 3, DATETIME('now'), DATETIME('now'), DATETIME('now'));
