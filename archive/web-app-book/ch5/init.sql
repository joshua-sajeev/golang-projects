-- Create tables for the blog application
CREATE TABLE IF NOT EXISTS pages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    page_title VARCHAR(255) NOT NULL,
    page_content TEXT NOT NULL,
    page_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    page_guid VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    page_id INT NOT NULL,
    comment_guid VARCHAR(255) NOT NULL UNIQUE,
    comment_name VARCHAR(255) NOT NULL,
    comment_email VARCHAR(255) NOT NULL,
    comment_text TEXT NOT NULL,
    comment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (page_id) REFERENCES pages(id)
);

-- Add indexes for better performance
CREATE INDEX idx_page_guid ON pages(page_guid);
CREATE INDEX idx_comment_page_id ON comments(page_id);
CREATE INDEX idx_comment_date ON comments(comment_date);
