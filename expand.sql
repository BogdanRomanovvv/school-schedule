-- Очищаем расписание
DELETE FROM schedule;

-- ── Новые предметы ────────────────────────────────────────────────────────────
INSERT INTO subjects(name) VALUES
    ('Обществознание'),
    ('ОБЖ'),
    ('Технология'),
    ('Музыка'),
    ('ИЗО')
ON CONFLICT DO NOTHING;

-- ── Новые учителя ─────────────────────────────────────────────────────────────
INSERT INTO teachers(name, max_hours_per_week) VALUES
    ('Ковалёв Д.С.',  36),
    ('Белова И.Н.',   24)
ON CONFLICT DO NOTHING;

-- ── Назначение предметов учителям ─────────────────────────────────────────────
INSERT INTO teacher_subjects(teacher_id, subject_id)
SELECT t.id, s.id FROM teachers t, subjects s
WHERE t.name='Петров В.С.' AND s.name='Обществознание'
ON CONFLICT DO NOTHING;

INSERT INTO teacher_subjects(teacher_id, subject_id)
SELECT t.id, s.id FROM teachers t, subjects s
WHERE t.name='Ковалёв Д.С.' AND s.name IN ('ОБЖ','Технология')
ON CONFLICT DO NOTHING;

INSERT INTO teacher_subjects(teacher_id, subject_id)
SELECT t.id, s.id FROM teachers t, subjects s
WHERE t.name='Белова И.Н.' AND s.name IN ('Музыка','ИЗО')
ON CONFLICT DO NOTHING;

-- ── 1-2 классы: +Технология 2, Музыка 1, ИЗО 1 ──────────────────────────────
INSERT INTO curriculum(class_id, subject_id, hours_per_week)
SELECT c.id, s.id, v.h FROM classes c, subjects s,
(VALUES ('Технология',2),('Музыка',1),('ИЗО',1)) AS v(sn,h)
WHERE c.name IN ('1','2') AND s.name=v.sn
ON CONFLICT DO NOTHING;

-- ── 3-4 классы: +Технология 2, Музыка 1, ИЗО 1 ──────────────────────────────
INSERT INTO curriculum(class_id, subject_id, hours_per_week)
SELECT c.id, s.id, v.h FROM classes c, subjects s,
(VALUES ('Технология',2),('Музыка',1),('ИЗО',1)) AS v(sn,h)
WHERE c.name IN ('3','4') AND s.name=v.sn
ON CONFLICT DO NOTHING;

-- ── 5-6 классы: +Общество 1, ОБЖ 1, Технология 2, Музыка 1, ИЗО 1 ──────────
INSERT INTO curriculum(class_id, subject_id, hours_per_week)
SELECT c.id, s.id, v.h FROM classes c, subjects s,
(VALUES ('Обществознание',1),('ОБЖ',1),('Технология',2),('Музыка',1),('ИЗО',1)) AS v(sn,h)
WHERE c.name IN ('5','6') AND s.name=v.sn
ON CONFLICT DO NOTHING;

-- ── 7-8 классы: +Общество 1, ОБЖ 1, Технология 2, Музыка 1, ИЗО 1 ──────────
INSERT INTO curriculum(class_id, subject_id, hours_per_week)
SELECT c.id, s.id, v.h FROM classes c, subjects s,
(VALUES ('Обществознание',1),('ОБЖ',1),('Технология',2),('Музыка',1),('ИЗО',1)) AS v(sn,h)
WHERE c.name IN ('7','8') AND s.name=v.sn
ON CONFLICT DO NOTHING;

-- ── 9 класс: +Обществознание 1, ОБЖ 1, Технология 1 ─────────────────────────
INSERT INTO curriculum(class_id, subject_id, hours_per_week)
SELECT c.id, s.id, v.h FROM classes c, subjects s,
(VALUES ('Обществознание',1),('ОБЖ',1),('Технология',1)) AS v(sn,h)
WHERE c.name='9' AND s.name=v.sn
ON CONFLICT DO NOTHING;

-- ── 10-11 классы: +Обществознание 2, ОБЖ 1 ──────────────────────────────────
INSERT INTO curriculum(class_id, subject_id, hours_per_week)
SELECT c.id, s.id, v.h FROM classes c, subjects s,
(VALUES ('Обществознание',2),('ОБЖ',1)) AS v(sn,h)
WHERE c.name IN ('10','11') AND s.name=v.sn
ON CONFLICT DO NOTHING;
