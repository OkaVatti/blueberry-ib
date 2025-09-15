-- apps/backend/migrations/0002_threads_and_content_html.sql
-- Threads table (each thread belongs to a board and a user)
CREATE TABLE
    IF NOT EXISTS threads (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        board_id uuid NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
        user_id uuid REFERENCES users (id) ON DELETE SET NULL,
        title text,
        created_at timestamptz DEFAULT now (),
        updated_at timestamptz DEFAULT now ()
    );

-- Add thread_id to posts so posts can belong to a thread (nullable)
ALTER TABLE posts
ADD COLUMN IF NOT EXISTS thread_id uuid REFERENCES threads (id) ON DELETE CASCADE;

-- Add content_html to posts and comments to store rendered & sanitized HTML
ALTER TABLE posts
ADD COLUMN IF NOT EXISTS content_html text;

ALTER TABLE comments
ADD COLUMN IF NOT EXISTS content_html text;