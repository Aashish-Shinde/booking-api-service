USE booking_api;

-- Delete all dummy data in reverse order of insertion
DELETE FROM bookings;
DELETE FROM availability_exceptions;
DELETE FROM availability;
DELETE FROM coaches;
DELETE FROM users;
