CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE roles (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name TEXT UNIQUE NOT NULL
);

INSERT INTO roles (name) VALUES
                             ('admin'),
                             ('node'),
                             ('phenophase')
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users ADD COLUMN role_id UUID;

UPDATE users
SET role_id = r.id
FROM roles r
WHERE r.name = CASE users.role
                   WHEN 'admin' THEN 'admin'
                   WHEN 'worker' THEN 'node'
                   WHEN 'agronom' THEN 'phenophase'
                   ELSE users.role
    END;

ALTER TABLE users
    ALTER COLUMN role_id SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT fk_users_role
        FOREIGN KEY (role_id) REFERENCES roles(id);

ALTER TABLE users
    DROP COLUMN role;

CREATE TABLE checklists (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            name TEXT NOT NULL,
                            code TEXT UNIQUE NOT NULL,
                            allowed_role_id UUID NOT NULL,
                            created_at TIMESTAMP DEFAULT now(),

                            CONSTRAINT fk_checklists_role
                                FOREIGN KEY (allowed_role_id) REFERENCES roles(id)
);

INSERT INTO checklists (name, code, allowed_role_id)
VALUES
    (
        'Мониторинг растворных узлов',
        'default',
        (SELECT id FROM roles WHERE name = 'node')
    ),
    (
        'Мониторинг фенофаз',
        'sort_control',
        (SELECT id FROM roles WHERE name = 'phenophase')
    )
ON CONFLICT (code) DO NOTHING;

ALTER TABLE questions
    ADD COLUMN checklist_id UUID,
    ADD COLUMN formula TEXT;

UPDATE questions
SET checklist_id = (
    SELECT id FROM checklists WHERE code = 'default'
)
WHERE checklist_id IS NULL;

ALTER TABLE questions
    ALTER COLUMN checklist_id SET NOT NULL;

ALTER TABLE questions
    ADD CONSTRAINT fk_questions_checklist
        FOREIGN KEY (checklist_id) REFERENCES checklists(id);

ALTER TABLE reports
    ADD COLUMN checklist_id UUID,
    ADD COLUMN metadata JSONB;

UPDATE reports
SET checklist_id = (
    SELECT id FROM checklists WHERE code = 'default'
)
WHERE checklist_id IS NULL;

ALTER TABLE reports
    ALTER COLUMN checklist_id SET NOT NULL;

ALTER TABLE reports
    ADD CONSTRAINT fk_reports_checklist
        FOREIGN KEY (checklist_id) REFERENCES checklists(id);

ALTER TABLE reports
    ALTER COLUMN place DROP NOT NULL;

ALTER TABLE answers
    ADD COLUMN result TEXT;