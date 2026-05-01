CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE phenophases
    ADD COLUMN IF NOT EXISTS min_critical_temperature NUMERIC,
    ADD COLUMN IF NOT EXISTS critical_temperature NUMERIC;

ALTER TABLE questions
    ADD COLUMN IF NOT EXISTS technical_code TEXT;

UPDATE questions q
SET technical_code = 'actual_temperature'
    FROM checklists c
WHERE q.checklist_id = c.id
  AND c.code = 'sort_control'
  AND q.technical_code IS NULL
  AND q.order_index = 2;

UPDATE questions q
SET technical_code = 'min_critical_temperature'
    FROM checklists c
WHERE q.checklist_id = c.id
  AND c.code = 'sort_control'
  AND q.technical_code IS NULL
  AND q.order_index = 3;

UPDATE questions q
SET technical_code = 'critical_temperature'
    FROM checklists c
WHERE q.checklist_id = c.id
  AND c.code = 'sort_control'
  AND q.technical_code IS NULL
  AND q.order_index = 4;