USE booking_api;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS availability_exceptions;
DROP TABLE IF EXISTS availability;
DROP TABLE IF EXISTS coaches;
DROP TABLE IF EXISTS users;
