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
                                              formula TEXT NOT NULL,
                                              created_at TIMESTAMP DEFAULT now(),
                                              updated_at TIMESTAMP DEFAULT now(),
                                              UNIQUE(question_id, phenophase_id)
);

INSERT INTO varieties (name, description, priority)
VALUES
    ('Gala', 'Популярный сладкий сорт яблони с красно-жёлтой окраской.', 'medium'),
    ('Fuji', 'Сладкий и хрустящий сорт яблони японского происхождения.', 'high'),
    ('Granny Smith', 'Зелёный кислый сорт яблони, хорошо известный для свежего употребления и кулинарии.', 'medium'),
    ('Golden Delicious', 'Жёлтый сладкий сорт яблони универсального назначения.', 'low'),
    ('Honeycrisp', 'Сочный и хрустящий сорт яблони с выраженной сладостью.', 'high')
ON CONFLICT (name) DO NOTHING;

INSERT INTO phenophases (name, description, order_index)
VALUES
    ('Период покоя', 'Зимний покой растения до начала активной вегетации.', 1),
    ('Набухание почек', 'Почки увеличиваются в размере перед распусканием.', 2),
    ('Распускание почек', 'Начало раскрытия почек и появления зелёной ткани.', 3),
    ('Зелёный конус', 'Появление зелёных кончиков листьев из почек.', 4),
    ('Розовый бутон', 'Бутоны сформированы, но цветки ещё не раскрыты.', 5),
    ('Цветение', 'Открытие цветков и период опыления.', 6),
    ('Опадение лепестков', 'Завершение цветения, лепестки начинают опадать.', 7),
    ('Завязь плодов', 'Формирование начальной завязи после опыления.', 8),
    ('Рост плодов', 'Активное увеличение размера плодов.', 9),
    ('Созревание', 'Окрашивание и достижение зрелости плодов.', 10)
ON CONFLICT (name) DO NOTHING;

WITH phenophase_checklist AS (
    SELECT id FROM checklists WHERE code = 'sort_control'
),
     inserted_questions AS (
         INSERT INTO questions (text, order_index, is_active, checklist_id, formula, image_url)
             VALUES
                 ('Оценка риска', 1, true, (SELECT id FROM phenophase_checklist), NULL, NULL),
                 ('T° факт, °C (min)', 2, true, (SELECT id FROM phenophase_checklist), NULL, NULL),
                 ('Мин. крит. порог, °C', 3, true, (SELECT id FROM phenophase_checklist), NULL, NULL),
                 ('Крит. температурный порог, °C', 4, true, (SELECT id FROM phenophase_checklist), NULL, NULL)
             RETURNING id, text
     )
INSERT INTO question_phenophase_formulas (question_id, phenophase_id, formula)
SELECT
    q.id,
    p.id,
    CASE
        WHEN q.text = 'Оценка риска' THEN '=Норма'
        WHEN q.text = 'T° факт, °C (min)' THEN
            CASE
                WHEN p.order_index <= 3 THEN '>=-10'
                WHEN p.order_index <= 6 THEN '>=-3'
                ELSE '>=0'
                END
        WHEN q.text = 'Мин. крит. порог, °C' THEN
            CASE
                WHEN p.order_index <= 3 THEN '>=-15'
                WHEN p.order_index <= 6 THEN '>=-5'
                ELSE '>=-2'
                END
        WHEN q.text = 'Крит. температурный порог, °C' THEN
            CASE
                WHEN p.order_index <= 3 THEN '>=-20'
                WHEN p.order_index <= 6 THEN '>=-7'
                ELSE '>=-4'
                END
        END
FROM inserted_questions q
         CROSS JOIN phenophases p
ON CONFLICT (question_id, phenophase_id) DO NOTHING;