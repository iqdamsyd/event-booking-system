ALTER TABLE bookings ADD COLUMN status_new booking_statuses;

UPDATE bookings SET status_new = status::text::booking_statuses;

ALTER TABLE bookings DROP COLUMN status;

ALTER TABLE bookings RENAME COLUMN status_new TO status;