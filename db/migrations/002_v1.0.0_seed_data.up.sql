USE booking_api;

-- Seed dummy users
INSERT INTO users (name, timezone) VALUES
('John Doe', 'America/New_York'),
('Jane Smith', 'Europe/London'),
('Mike Johnson', 'Asia/Tokyo'),
('Sarah Williams', 'Australia/Sydney'),
('David Brown', 'America/Los_Angeles');

-- Seed dummy coaches
INSERT INTO coaches (name, timezone) VALUES
('Coach Alice', 'America/New_York'),
('Coach Bob', 'Europe/London'),
('Coach Charlie', 'Asia/Tokyo');

-- Seed availability for coaches (Monday to Friday)
-- Coach Alice: Mon-Fri 9:00 AM - 5:00 PM
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(1, 1, '09:00:00', '17:00:00'), -- Monday
(1, 2, '09:00:00', '17:00:00'), -- Tuesday
(1, 3, '09:00:00', '17:00:00'), -- Wednesday
(1, 4, '09:00:00', '17:00:00'), -- Thursday
(1, 5, '09:00:00', '17:00:00'), -- Friday
-- Coach Bob: Tue-Sat 10:00 AM - 6:00 PM
(2, 2, '10:00:00', '18:00:00'), -- Tuesday
(2, 3, '10:00:00', '18:00:00'), -- Wednesday
(2, 4, '10:00:00', '18:00:00'), -- Thursday
(2, 5, '10:00:00', '18:00:00'), -- Friday
(2, 6, '10:00:00', '18:00:00'), -- Saturday
-- Coach Charlie: Mon, Wed, Fri 8:00 AM - 4:00 PM
(3, 1, '08:00:00', '16:00:00'), -- Monday
(3, 3, '08:00:00', '16:00:00'), -- Wednesday
(3, 5, '08:00:00', '16:00:00'); -- Friday

-- Seed availability exceptions (e.g., holidays or days off)
INSERT INTO availability_exceptions (coach_id, date, is_available) VALUES
(1, DATE_ADD(CURDATE(), INTERVAL 7 DAY), FALSE),  -- Coach Alice off next week
(2, DATE_ADD(CURDATE(), INTERVAL 14 DAY), FALSE); -- Coach Bob off two weeks from now

-- Seed some bookings
INSERT INTO bookings (user_id, coach_id, start_time, end_time, status, idempotency_key) VALUES
(1, 1, '2026-04-06 09:00:00', '2026-04-06 10:00:00', 'ACTIVE', 'idem-001'),
(2, 2, '2026-04-07 10:00:00', '2026-04-07 11:00:00', 'ACTIVE', 'idem-002'),
(3, 3, '2026-04-08 08:00:00', '2026-04-08 09:00:00', 'ACTIVE', 'idem-003'),
(4, 1, '2026-04-09 14:00:00', '2026-04-09 15:00:00', 'ACTIVE', 'idem-004'),
(5, 2, '2026-04-10 15:00:00', '2026-04-10 16:00:00', 'COMPLETED', 'idem-005');
