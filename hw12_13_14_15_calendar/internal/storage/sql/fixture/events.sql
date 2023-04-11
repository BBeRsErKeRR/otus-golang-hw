TRUNCATE TABLE events;
INSERT INTO events(id, title, date, end_date, description, user_id, remind_date) 
VALUES 
    ('2', 'title 2', (NOW() - INTERVAL '1 day'), (NOW() - INTERVAL '1 day'), '', '1', '0001-01-01 00:00:00'),
    ('1', 'title 1', NOW(), (NOW() + INTERVAL '4 hours'), '', '1', '0001-01-01 00:00:00'),
    ('3', 'title 3', (NOW()+ INTERVAL '1 day'), (NOW() + INTERVAL '26 hours'), '', '1', '0001-01-01 00:00:00'),
    ('4', 'title 4', (NOW()+ INTERVAL '14 days'), (NOW() + INTERVAL '15 days'), '', '1', '0001-01-01 00:00:00');