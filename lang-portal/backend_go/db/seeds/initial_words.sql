-- Clear existing data
DELETE FROM word_review_items;
DELETE FROM study_sessions;
DELETE FROM words_groups;
DELETE FROM words;
DELETE FROM groups;
DELETE FROM study_activities;

-- Insert groups
INSERT INTO groups (id, name, description) VALUES
(1, 'Greetings', 'Essential Japanese greetings and phrases'),
(2, 'Numbers', 'Basic Japanese numbers'),
(3, 'Daily Life', 'Common everyday words'),
(4, 'Family', 'Family-related vocabulary'),
(5, 'Food', 'Common food and drink terms');

-- Insert words
INSERT INTO words (id, original, translated, pronunciation, example_sentence, example_translation) VALUES
-- Greetings
(1, 'Hello', 'こんにちは', 'Konnichiwa', 'こんにちは、お元気ですか？', 'Hello, how are you?'),
(2, 'Good morning', 'おはようございます', 'Ohayou gozaimasu', 'おはようございます、今日もいい天気ですね。', 'Good morning, nice weather today.'),
(3, 'Good evening', 'こんばんは', 'Konbanwa', 'こんばんは、お疲れ様でした。', 'Good evening, good work today.'),
(4, 'Thank you', 'ありがとうございます', 'Arigatou gozaimasu', 'ご親切にありがとうございます。', 'Thank you for your kindness.'),

-- Numbers
(5, 'One', '一', 'Ichi', '一番好きな食べ物は寿司です。', 'My favorite food is sushi.'),
(6, 'Two', '二', 'Ni', '二人で映画を見に行きました。', 'We went to see a movie together.'),
(7, 'Three', '三', 'San', '三日後に日本に行きます。', 'I will go to Japan in three days.'),
(8, 'Four', '四', 'Shi/Yon', '四季があります。', 'There are four seasons.'),

-- Daily Life
(9, 'Water', '水', 'Mizu', '水を飲みたいです。', 'I want to drink water.'),
(10, 'Book', '本', 'Hon', 'この本は面白いです。', 'This book is interesting.'),
(11, 'Train', '電車', 'Densha', '電車で学校に行きます。', 'I go to school by train.'),
(12, 'Time', '時間', 'Jikan', '時間がありません。', 'There is no time.'),

-- Family
(13, 'Mother', '母', 'Haha', '母は料理が上手です。', 'My mother is good at cooking.'),
(14, 'Father', '父', 'Chichi', '父は会社員です。', 'My father is a company employee.'),
(15, 'Sister', '姉', 'Ane', '姉は医者です。', 'My older sister is a doctor.'),
(16, 'Brother', '兄', 'Ani', '兄は東京に住んでいます。', 'My older brother lives in Tokyo.'),

-- Food
(17, 'Rice', 'ご飯', 'Gohan', 'ご飯を食べます。', 'I eat rice.'),
(18, 'Fish', '魚', 'Sakana', '魚が好きです。', 'I like fish.'),
(19, 'Tea', 'お茶', 'Ocha', 'お茶を飲みませんか？', 'Would you like to drink tea?'),
(20, 'Meat', '肉', 'Niku', '肉を食べたいです。', 'I want to eat meat.');

-- Link words to groups
INSERT INTO words_groups (word_id, group_id) VALUES
-- Greetings
(1, 1), (2, 1), (3, 1), (4, 1),
-- Numbers
(5, 2), (6, 2), (7, 2), (8, 2),
-- Daily Life
(9, 3), (10, 3), (11, 3), (12, 3),
-- Family
(13, 4), (14, 4), (15, 4), (16, 4),
-- Food
(17, 5), (18, 5), (19, 5), (20, 5);

-- Insert study activities
INSERT INTO study_activities (id, name, description) VALUES
(1, 'Vocabulary Review', '語彙の復習'),
(2, 'Pronunciation Practice', '発音練習'),
(3, 'Kanji Study', '漢字学習'),
(4, 'Sentence Building', '文章作成');
