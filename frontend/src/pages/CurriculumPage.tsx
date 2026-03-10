import { useEffect, useState } from 'react';
import * as api from '../api';
import type { Class, Subject, Curriculum } from '../types';

export default function CurriculumPage() {
    const [classes, setClasses] = useState<Class[]>([]);
    const [subjects, setSubjects] = useState<Subject[]>([]);
    const [curriculum, setCurriculum] = useState<Curriculum[]>([]);
    const [classId, setClassId] = useState<number>(0);
    const [subjectId, setSubjectId] = useState<number>(0);
    const [hours, setHours] = useState<number>(2);
    const [filterClass, setFilterClass] = useState<number>(0);
    const [error, setError] = useState('');

    const load = async () => {
        try {
            const [cls, sub, cur] = await Promise.all([api.getClasses(), api.getSubjects(), api.getCurriculum(filterClass || undefined)]);
            setClasses(cls); setSubjects(sub); setCurriculum(cur);
        } catch (e: unknown) { setError(String(e)); }
    };

    useEffect(() => { load(); }, [filterClass]);

    const handleSave = async () => {
        if (!classId || !subjectId || hours < 1) return;
        try {
            await api.upsertCurriculum({ class_id: classId, subject_id: subjectId, hours_per_week: hours });
            load();
        } catch (e: unknown) { setError(String(e)); }
    };

    const handleDelete = async (classId: number, subjectId: number) => {
        if (!confirm('Удалить запись из учебного плана?')) return;
        try { await api.deleteCurriculum(classId, subjectId); load(); }
        catch (e: unknown) { setError(String(e)); }
    };

    const getClassName = (id: number) => classes.find(c => c.id === id)?.name ?? id;
    const getSubjectName = (id: number) => subjects.find(s => s.id === id)?.name ?? id;

    return (
        <div>
            <div className="card">
                <h2>Учебный план</h2>
                {error && <div className="alert alert-error">{error}</div>}
                <div className="form-row">
                    <div>
                        <label>Класс</label>
                        <select value={classId} onChange={e => setClassId(Number(e.target.value))}>
                            <option value={0}>-- выберите --</option>
                            {classes.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                        </select>
                    </div>
                    <div>
                        <label>Предмет</label>
                        <select value={subjectId} onChange={e => setSubjectId(Number(e.target.value))}>
                            <option value={0}>-- выберите --</option>
                            {subjects.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                        </select>
                    </div>
                    <div>
                        <label>Часов/нед.</label>
                        <input type="number" value={hours} min={1} max={7}
                            onChange={e => setHours(Number(e.target.value))} style={{ width: 80 }} />
                    </div>
                    <button className="btn-primary" disabled={!classId || !subjectId} onClick={handleSave}>
                        Сохранить
                    </button>
                </div>

                <div className="form-row" style={{ marginBottom: 16 }}>
                    <div>
                        <label>Фильтр по классу</label>
                        <select value={filterClass} onChange={e => setFilterClass(Number(e.target.value))}>
                            <option value={0}>Все классы</option>
                            {classes.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                        </select>
                    </div>
                </div>

                <table>
                    <thead>
                        <tr><th>Класс</th><th>Предмет</th><th>Часов/нед.</th><th>Действия</th></tr>
                    </thead>
                    <tbody>
                        {curriculum.map(c => (
                            <tr key={`${c.class_id}-${c.subject_id}`}>
                                <td>{getClassName(c.class_id)}</td>
                                <td>{getSubjectName(c.subject_id)}</td>
                                <td>{c.hours_per_week}</td>
                                <td>
                                    <button className="btn-danger" onClick={() => handleDelete(c.class_id, c.subject_id)}>Удалить</button>
                                </td>
                            </tr>
                        ))}
                        {curriculum.length === 0 && (
                            <tr><td colSpan={4} style={{ textAlign: 'center', color: '#9ca3af' }}>Учебный план пуст</td></tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
