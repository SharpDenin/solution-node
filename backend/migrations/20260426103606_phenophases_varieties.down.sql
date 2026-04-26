DROP TABLE IF EXISTS question_phenophase_formulas;

ALTER TABLE questions
    DROP COLUMN IF EXISTS image_url;

ALTER TABLE reports
    DROP COLUMN IF EXISTS variety_id,
    DROP COLUMN IF EXISTS phenophase_id;

DROP TABLE IF EXISTS phenophases;

DROP TABLE IF EXISTS varieties;

ALTER TABLE users
    DROP COLUMN IF EXISTS position;