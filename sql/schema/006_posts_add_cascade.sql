-- +goose Up
ALTER TABLE posts DROP CONSTRAINT fk_feed;

ALTER TABLE posts
ADD CONSTRAINT fk_feed
FOREIGN KEY (feed_id)
REFERENCES feeds(id)
ON DELETE CASCADE;


-- +goose Down
ALTER TABLE posts DROP CONSTRAINT fk_feed;

ALTER TABLE posts
ADD CONSTRAINT fk_feed
FOREIGN KEY (feed_id)
REFERENCES feeds(id);
