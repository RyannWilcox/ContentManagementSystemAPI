CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  author VARCHAR(100),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE post_media (
  post_id BIGINT,
  media_id BIGINT,

  PRIMARY KEY(post_id, media_id),

  CONSTRAINT `fk_post_media_post` 
    FOREIGN KEY (`post_id`) 
    REFERENCES `posts`(`id`) 
    ON DELETE CASCADE,

  CONSTRAINT `fk_post_media_media` 
    FOREIGN KEY (`media_id`) 
    REFERENCES `media`(`id`) 
    ON DELETE CASCADE
);
