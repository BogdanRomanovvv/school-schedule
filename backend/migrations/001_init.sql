-- 001_init.sql
-- Полная схема базы данных для системы генерации расписания

CREATE TABLE IF NOT EXISTS classes (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS subjects (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS teachers (
    id                SERIAL PRIMARY KEY,
    name              VARCHAR(150) NOT NULL,
    max_hours_per_week INT NOT NULL DEFAULT 36,
    -- Если задан, учитель является классным руководителем и ведёт уроки ТОЛЬКО в этом классе
    homeroom_class_id INT REFERENCES classes(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS teacher_subjects (
    teacher_id INT NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    subject_id INT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    PRIMARY KEY (teacher_id, subject_id)
);

CREATE TABLE IF NOT EXISTS curriculum (
    id            SERIAL PRIMARY KEY,
    class_id      INT NOT NULL REFERENCES classes(id)  ON DELETE CASCADE,
    subject_id    INT NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    hours_per_week INT NOT NULL CHECK (hours_per_week > 0 AND hours_per_week <= 7),
    UNIQUE (class_id, subject_id)
);

CREATE TABLE IF NOT EXISTS schedule (
    id            SERIAL PRIMARY KEY,
    class_id      INT NOT NULL REFERENCES classes(id)   ON DELETE CASCADE,
    subject_id    INT NOT NULL REFERENCES subjects(id)  ON DELETE CASCADE,
    teacher_id    INT NOT NULL REFERENCES teachers(id)  ON DELETE CASCADE,
    day           SMALLINT NOT NULL CHECK (day   BETWEEN 0 AND 4),
    lesson_number SMALLINT NOT NULL CHECK (lesson_number BETWEEN 0 AND 6),
    UNIQUE (class_id,   day, lesson_number),
    UNIQUE (teacher_id, day, lesson_number)
);

-- ─── Тестовые данные (можно удалить в продакшн) ───────────────────────────────

INSERT INTO classes(name) VALUES ('5А'),('6А'),('7А'),('8А'),('9А'),('10А'),('11А')
ON CONFLICT DO NOTHING;

INSERT INTO subjects(name) VALUES
    ('Математика'),('Физика'),('Химия'),('Биология'),('История'),
    ('Русский язык'),('Литература'),('Английский язык'),('Информатика'),('Физкультура')
ON CONFLICT DO NOTHING;

INSERT INTO teachers(name, max_hours_per_week) VALUES
    ('Иванова А.П.',  30),
    ('Петров В.С.',   28),
    ('Сидорова Н.К.', 32),
    ('Козлов М.И.',   25),
    ('Новикова О.Э.', 30)
ON CONFLICT DO NOTHING;

-- Назначение предметов учителям
INSERT INTO teacher_subjects(teacher_id, subject_id)
SELECT t.id, s.id FROM teachers t, subjects s
WHERE (t.name='Иванова А.П.'  AND s.name IN ('Математика','Физика'))
   OR (t.name='Петров В.С.'   AND s.name IN ('История','Обществознание'))
   OR (t.name='Сидорова Н.К.' AND s.name IN ('Русский язык','Литература'))
   OR (t.name='Козлов М.И.'   AND s.name IN ('Информатика','Математика'))
   OR (t.name='Новикова О.Э.' AND s.name IN ('Биология','Химия','Физкультура','Английский язык'))
ON CONFLICT DO NOTHING;
