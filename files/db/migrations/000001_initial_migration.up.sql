
CREATE TABLE IF NOT EXISTS posts (
	id BIGSERIAL PRIMARY KEY,
	title  VARCHAR(50) NOT NULL,
	content VARCHAR(1024) NOT NULL,
	status varchar(20) NOT NULL,
	publish_date timestamp
);

CREATE TABLE IF NOT EXISTS tags (
	id BIGSERIAL PRIMARY KEY,
	label varchar(20) NOT NULL UNIQUE

);
CREATE INDEX idx_label ON tags(label);

create TABLE IF NOT EXISTS post_tags (
	id BIGSERIAL PRIMARY KEY,
	tag_id BIGINT NOT NULL,
	post_id BIGINT NOT NULL,

	FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
	FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE


);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_post_tags ON post_tags(tag_id,post_id);

CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	email VARCHAR(20) UNIQUE,
	password VARCHAR(255) NOT NULL,
	role varchar(20) NOT NULL
);
