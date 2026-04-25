ALTER TABLE answers
    DROP COLUMN IF EXISTS result;

ALTER TABLE reports
    DROP CONSTRAINT IF EXISTS fk_reports_checklist;

ALTER TABLE reports
    DROP COLUMN IF EXISTS checklist_id,
    DROP COLUMN IF EXISTS metadata;

ALTER TABLE reports
    ALTER COLUMN place SET NOT NULL;

ALTER TABLE questions
    DROP CONSTRAINT IF EXISTS fk_questions_checklist;

ALTER TABLE questions
    DROP COLUMN IF EXISTS checklist_id,
    DROP COLUMN IF EXISTS formula;

DROP TABLE IF EXISTS checklists;

ALTER TABLE users ADD COLUMN role TEXT;

UPDATE users
SET role = r.name
FROM roles r
WHERE users.role_id = r.id;

UPDATE users
SET role = CASE role
               WHEN 'node' THEN 'worker'
               WHEN 'phenophase' THEN 'agronom'
               ELSE role
    END;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS fk_users_role;

ALTER TABLE users
    DROP COLUMN IF EXISTS role_id;

ALTER TABLE users
    ALTER COLUMN role SET NOT NULL;

DROP TABLE IF EXISTS roles;