-- TODO: Create posts table
-- - Add primary key id column as SERIAL
-- - Add required title column with varchar(255)
-- - Add required content column as text
-- - Add optional author column
-- - Add timestamps for created_at and updated_at

-- TODO: Create post_media junction table
-- - Add composite primary key (post_id, media_id)
-- - Add foreign key constraint for post_id referencing posts table
-- - Add foreign key constraint for media_id referencing media table
-- - Add cascade delete for both foreign keys

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
