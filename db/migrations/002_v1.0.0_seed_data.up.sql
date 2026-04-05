USE booking_api;

-- ============================================================================
-- USERS (Multiple Users for testing different scenarios)
-- ============================================================================
INSERT INTO users (name, timezone) VALUES
-- Base users
('John Doe', 'America/New_York'),           -- ID: 1
('Jane Smith', 'Europe/London'),            -- ID: 2
('Mike Johnson', 'Asia/Tokyo'),             -- ID: 3
('Sarah Williams', 'Australia/Sydney'),     -- ID: 4
('David Brown', 'America/Los_Angeles'),     -- ID: 5
-- Additional users for scenario testing
('Emma Davis', 'America/Chicago'),          -- ID: 6
('Alex Martinez', 'Europe/Paris'),          -- ID: 7
('Lisa Anderson', 'Asia/Singapore'),        -- ID: 8
('Tom Wilson', 'Canada/Toronto'),           -- ID: 9
('Rachel Green', 'America/Denver');         -- ID: 10

-- ============================================================================
-- COACHES (Multiple Coaches with different availability patterns)
-- ============================================================================
INSERT INTO coaches (name, timezone) VALUES
-- Base coaches
('Coach Alice', 'America/New_York'),        -- ID: 1 - Mon-Fri standard hours
('Coach Bob', 'Europe/London'),             -- ID: 2 - Tue-Sat extended hours
('Coach Charlie', 'Asia/Tokyo'),            -- ID: 3 - Mon, Wed, Fri only
-- Additional coaches for variety
('Coach Diana', 'America/Los_Angeles'),     -- ID: 4 - Mon-Thu, early hours
('Coach Edward', 'Europe/Berlin'),          -- ID: 5 - Tue-Sat, flexible slots
('Coach Fiona', 'Australia/Melbourne'),     -- ID: 6 - All days available
('Coach George', 'America/Chicago'),        -- ID: 7 - Limited availability (Wed-Thu)
('Coach Hannah', 'Asia/Singapore');         -- ID: 8 - Daily with multiple slots

-- ============================================================================
-- AVAILABILITY (Weekly recurring availability for coaches)
-- ============================================================================

-- Coach Alice (ID: 1): Mon-Fri 9:00 AM - 5:00 PM (New York time)
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(1, 1, '09:00', '17:00'), -- Monday
(1, 2, '09:00', '17:00'), -- Tuesday
(1, 3, '09:00', '17:00'), -- Wednesday
(1, 4, '09:00', '17:00'), -- Thursday
(1, 5, '09:00', '17:00'); -- Friday

-- Coach Bob (ID: 2): Tue-Sat 10:00 AM - 6:00 PM (London time) with lunch break
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(2, 2, '10:00', '12:30'), -- Tuesday morning
(2, 2, '13:30', '18:00'), -- Tuesday afternoon
(2, 3, '10:00', '18:00'), -- Wednesday full day
(2, 4, '10:00', '18:00'), -- Thursday full day
(2, 5, '10:00', '18:00'), -- Friday full day
(2, 6, '10:00', '18:00'); -- Saturday full day

-- Coach Charlie (ID: 3): Mon, Wed, Fri 8:00 AM - 4:00 PM (Tokyo time)
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(3, 1, '08:00', '16:00'), -- Monday
(3, 3, '08:00', '16:00'), -- Wednesday
(3, 5, '08:00', '16:00'); -- Friday

-- Coach Diana (ID: 4): Mon-Thu 7:00 AM - 3:00 PM (LA time)
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(4, 1, '07:00', '15:00'), -- Monday
(4, 2, '07:00', '15:00'), -- Tuesday
(4, 3, '07:00', '15:00'), -- Wednesday
(4, 4, '07:00', '15:00'); -- Thursday

-- Coach Edward (ID: 5): Tue-Sat with flexible slots (Berlin time)
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(5, 2, '09:00', '12:00'), -- Tuesday morning
(5, 2, '14:00', '17:00'), -- Tuesday afternoon
(5, 3, '10:00', '18:00'), -- Wednesday
(5, 4, '09:00', '12:00'), -- Thursday morning
(5, 4, '15:00', '19:00'), -- Thursday evening
(5, 5, '10:00', '18:00'), -- Friday
(5, 6, '11:00', '17:00'); -- Saturday

-- Coach Fiona (ID: 6): Every day with flexible hours (Melbourne time)
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(6, 0, '10:00', '18:00'), -- Sunday
(6, 1, '08:00', '18:00'), -- Monday
(6, 2, '08:00', '18:00'), -- Tuesday
(6, 3, '08:00', '18:00'), -- Wednesday
(6, 4, '08:00', '18:00'), -- Thursday
(6, 5, '08:00', '18:00'), -- Friday
(6, 6, '10:00', '18:00'); -- Saturday

-- Coach George (ID: 7): Limited availability (Chicago time)
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(7, 3, '14:00', '18:00'), -- Wednesday
(7, 4, '14:00', '18:00'); -- Thursday

-- Coach Hannah (ID: 8): Daily with multiple slots (Singapore time)
INSERT INTO availability (coach_id, day_of_week, start_time, end_time) VALUES
(8, 1, '08:00', '11:00'), -- Monday morning
(8, 1, '14:00', '17:00'), -- Monday afternoon
(8, 2, '08:00', '11:00'), -- Tuesday morning
(8, 2, '14:00', '17:00'), -- Tuesday afternoon
(8, 3, '08:00', '17:00'), -- Wednesday full day
(8, 4, '08:00', '11:00'), -- Thursday morning
(8, 4, '14:00', '17:00'), -- Thursday afternoon
(8, 5, '08:00', '17:00'), -- Friday full day
(8, 6, '09:00', '15:00'); -- Saturday

-- ============================================================================
-- AVAILABILITY EXCEPTIONS (One-time overrides for holidays/days off)
-- ============================================================================
INSERT INTO availability_exceptions (coach_id, date, is_available, start_time, end_time) VALUES
-- Days off
(1, DATE_ADD(CURDATE(), INTERVAL 7 DAY), FALSE, NULL, NULL),    -- Coach Alice off next week
(2, DATE_ADD(CURDATE(), INTERVAL 14 DAY), FALSE, NULL, NULL),   -- Coach Bob off two weeks
(3, DATE_ADD(CURDATE(), INTERVAL 3 DAY), FALSE, NULL, NULL),    -- Coach Charlie off in 3 days
-- Special availability (e.g., weekend availability)
(4, DATE_ADD(CURDATE(), INTERVAL 6 DAY), TRUE, '10:00', '14:00'), -- Coach Diana available Saturday
(5, DATE_ADD(CURDATE(), INTERVAL 1 DAY), TRUE, '19:00', '21:00'), -- Coach Edward evening session
-- Partial day off
(6, DATE_ADD(CURDATE(), INTERVAL 5 DAY), TRUE, '10:00', '13:00'); -- Coach Fiona available only morning

-- ============================================================================
-- BOOKINGS (Test data for various scenarios)
-- ============================================================================
INSERT INTO bookings (user_id, coach_id, start_time, end_time, status, idempotency_key) VALUES
-- Active bookings (various coaches and users)
(1, 1, '2026-04-06 09:00:00', '2026-04-06 10:00:00', 'ACTIVE', 'idem-booking-001'),
(2, 2, '2026-04-07 10:00:00', '2026-04-07 11:00:00', 'ACTIVE', 'idem-booking-002'),
(3, 3, '2026-04-08 08:00:00', '2026-04-08 09:00:00', 'ACTIVE', 'idem-booking-003'),
(4, 1, '2026-04-09 14:00:00', '2026-04-09 15:00:00', 'ACTIVE', 'idem-booking-004'),
(5, 2, '2026-04-10 15:00:00', '2026-04-10 16:00:00', 'ACTIVE', 'idem-booking-005'),
-- Additional active bookings with other coaches
(6, 4, '2026-04-07 11:00:00', '2026-04-07 12:00:00', 'ACTIVE', 'idem-booking-006'),
(7, 5, '2026-04-08 14:30:00', '2026-04-08 15:30:00', 'ACTIVE', 'idem-booking-007'),
(8, 6, '2026-04-06 14:00:00', '2026-04-06 15:00:00', 'ACTIVE', 'idem-booking-008'),
(9, 8, '2026-04-07 08:00:00', '2026-04-07 09:00:00', 'ACTIVE', 'idem-booking-009'),
(10, 1, '2026-04-11 10:00:00', '2026-04-11 11:00:00', 'ACTIVE', 'idem-booking-010'),
-- Completed bookings (test status filtering)
(1, 2, '2026-03-30 10:30:00', '2026-03-30 11:30:00', 'COMPLETED', 'idem-booking-011'),
(2, 1, '2026-03-31 09:00:00', '2026-03-31 10:00:00', 'COMPLETED', 'idem-booking-012'),
(3, 4, '2026-03-25 13:00:00', '2026-03-25 14:00:00', 'COMPLETED', 'idem-booking-013'),
(4, 6, '2026-03-20 16:00:00', '2026-03-20 17:00:00', 'COMPLETED', 'idem-booking-014'),
-- Cancelled bookings (test status filtering)
(5, 3, '2026-04-12 08:30:00', '2026-04-12 09:30:00', 'CANCELLED', 'idem-booking-015'),
(6, 5, '2026-04-13 15:00:00', '2026-04-13 16:00:00', 'CANCELLED', 'idem-booking-016'),
(7, 2, '2026-04-14 11:00:00', '2026-04-14 12:00:00', 'CANCELLED', 'idem-booking-017');
