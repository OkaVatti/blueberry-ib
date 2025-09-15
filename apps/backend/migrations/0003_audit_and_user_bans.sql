-- apps/backend/migrations/0003_audit_and_user_bans.sql
-- Add basic ban/suspension columns to users
ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_banned boolean DEFAULT false,
ADD COLUMN IF NOT EXISTS banned_until timestamptz;

-- Audit logs to record admin/moderation actions
CREATE TABLE
    IF NOT EXISTS audit_logs (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        actor_id uuid, -- user who performed the action (nullable e.g. system)
        action text NOT NULL, -- e.g. "ban_user", "delete_post"
        target_type text, -- "user","post","comment","thread","board"
        target_id uuid,
        details jsonb,
        created_at timestamptz DEFAULT now ()
    );

CREATE INDEX IF NOT EXISTS idx_audit_actor ON audit_logs (actor_id);

CREATE INDEX IF NOT EXISTS idx_audit_target ON audit_logs (target_type, target_id);