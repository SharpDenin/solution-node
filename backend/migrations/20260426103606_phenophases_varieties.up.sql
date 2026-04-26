CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE users
    ADD COLUMN position TEXT;

CREATE TABLE varieties (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           name TEXT UNIQUE NOT NULL,
                           description TEXT,
                           priority TEXT NOT NULL DEFAULT 'medium',
                           image_url TEXT,
                           created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE phenophases (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             name TEXT UNIQUE NOT NULL,
                             description TEXT,
                             image_url TEXT,
                             order_index INT NOT NULL,
                             created_at TIMESTAMP DEFAULT now()
);

ALTER TABLE reports
    ADD COLUMN variety_id UUID REFERENCES varieties(id),
    ADD COLUMN phenophase_id UUID REFERENCES phenophases(id);

ALTER TABLE questions
    ADD COLUMN image_url TEXT;

CREATE TABLE question_phenophase_formulas (
                                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                              question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
                                              phenophase_id UUID NOT NULL REFERENCES phenophases(id) ON DELETE CASCADE,
                                              formula TEXT,
                                              created_at TIMESTAMP DEFAULT now(),
                                              updated_at TIMESTAMP DEFAULT now(),
                                              UNIQUE(question_id, phenophase_id)
);

WITH phenophase_checklist AS (
    SELECT id FROM checklists WHERE code = 'sort_control'
)
INSERT INTO questions (text, order_index, is_active, checklist_id, formula, image_url)
SELECT
    v.text,
    v.order_index,
    true,
    pc.id,
    NULL,
    NULL
FROM phenophase_checklist pc
         CROSS JOIN (
    VALUES
        ('Оценка риска', 1),
        ('T° факт, °C (min)', 2),
        ('Мин. крит. порог, °C', 3),
        ('Крит. температурный порог, °C', 4)
) AS v(text, order_index);