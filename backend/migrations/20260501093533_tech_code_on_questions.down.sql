UPDATE questions
SET technical_code = NULL
WHERE technical_code IN (
                         'actual_temperature',
                         'min_critical_temperature',
                         'critical_temperature'
    );

ALTER TABLE questions
DROP COLUMN IF EXISTS technical_code;

ALTER TABLE phenophases
DROP COLUMN IF EXISTS min_critical_temperature,
    DROP COLUMN IF EXISTS critical_temperature;