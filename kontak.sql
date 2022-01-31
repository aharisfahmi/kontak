-- CREATE DATABASE `kontak`;

CREATE TABLE kontak.kontak (
	id varchar(128) NOT NULL PRIMARY KEY,
	name varchar(255),
	gender varchar(6),
	phone varchar(32),
	email varchar(255),
	created_at datetime DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

);

INSERT INTO kontak.kontak
(id, name, gender, phone, email, created_at, updated_at)
VALUES('1a5071bd-2960-4829-8adc-593e216b2de5', 'fulan', 'male', '628123456789', 'fulan@email.com', current_timestamp(), current_timestamp());
