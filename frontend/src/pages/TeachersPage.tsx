import { useEffect, useState } from 'react';
import * as api from '../api';
import type { Teacher, Subject } from '../types';

export default function TeachersPage() {
    const [teachers, setTeachers] = useState<Teacher[]>([]);
    const [subjects, setSubjects] = useState<Subject[]>([]);
    const [name, setName] = useState('');
    const [maxHours, setMaxHours] = useState(36);
    const [editing, setEditing] = useState<Teacher | null>(null);
    const [selectedTeacher, setSelectedTeacher] = useState<Teacher | null>(null);
    const [teacherSubjects, setTeacherSubjects] = useState<Subject[]>([]);
    const [assignSubjectId, setAssignSubjectId] = useState<number>(0);
    const [error, setError] = useState('');

    const load = async () => {
        try {
            const [t, s] = await Promise.all([api.getTeachers(), api.getSubjects()]);
            setTeachers(t); setSubjects(s);
        } catch (e: unknown) { setError(String(e)); }
    };

    useEffect(() => { load(); }, []);

    useEffect(() => {
        if (selectedTeacher) {
            api.getTeacherSubjects(selectedTeacher.id).then(setTeacherSubjects);
        }
    }, [selectedTeacher]);

    const handleSave = async () => {
        if (!name.trim()) return;
        try {
            if (editing) {
                await api.updateTeacher(editing.id, { name, max_hours_per_week: maxHours });
            } else {
                await api.createTeacher({ name, max_hours_per_week: maxHours });
            }
            setName(''); setMaxHours(36); setEditing(null); load();
        } catch (e: unknown) { setError(String(e)); }
    };

    const handleDelete = async (id: number) => {
        if (!confirm('Удалить учителя?')) return;
        try { await api.deleteTeacher(id); load(); }
        catch (e: unknown) { setError(String(e)); }
    };

    const handleAssign = async () => {
        if (!selectedTeacher || !assignSubjectId) return;
        try {
            await api.assignSubject(selectedTeacher.id, assignSubjectId);
            const subs = await api.getTeacherSubjects(selectedTeacher.id);
            setTeacherSubjects(subs);
        } catch (e: unknown) { setError(String(e)); }
    };

    const handleRemoveSubject = async (subjectId: number) => {
        if (!selectedTeacher) return;
        try {
            await api.removeSubject(selectedTeacher.id, subjectId);
            const subs = await api.getTeacherSubjects(selectedTeacher.id);
            setTeacherSubjects(subs);
        } catch (e: unknown) { setError(String(e)); }
    };

    const availableSubjects = subjects.filter(s => !teacherSubjects.find(ts => ts.id === s.id));

    return (
        <div>
            <div className="card">
                <h2>Учителя</h2>
                {error && <div className="alert alert-error">{error}</div>}
                <div className="form-row">
                    <div>
                        <label>ФИО</label>
                        <input value={name} onChange={e => setName(e.target.value)} placeholder="Иванова А.П." style={{ width: 220 }} />
                    </div>
                    <div>
                        <label>Макс. часов/нед.</label>
                        <input type="number" value={maxHours} min={1} max={40}
                            onChange={e => setMaxHours(Number(e.target.value))} style={{ width: 100 }} />
                    </div>
                    <button className="btn-primary" onClick={handleSave}>{editing ? 'Сохранить' : 'Добавить'}</button>
                    {editing && <button className="btn-secondary" onClick={() => { setEditing(null); setName(''); setMaxHours(36); }}>Отмена</button>}
                </div>

                <table>
                    <thead>
                        <tr><th>#</th><th>ФИО</th><th>Макс. часов/нед.</th><th>Действия</th></tr>
                    </thead>
                    <tbody>
                        {teachers.map((t, idx) => (
                            <tr key={t.id}>
                                <td>{idx + 1}</td>
                                <td>{t.name}</td>
                                <td>{t.max_hours_per_week}</td>
                                <td>
                                    <div className="action-btns">
                                        <button className="btn-secondary" onClick={() => { setEditing(t); setName(t.name); setMaxHours(t.max_hours_per_week); }}>Изменить</button>
                                        <button className="btn-success" onClick={() => setSelectedTeacher(t)}>Предметы</button>
                                        <button className="btn-danger" onClick={() => handleDelete(t.id)}>Удалить</button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {selectedTeacher && (
                <div className="card">
                    <h2>Предметы учителя: {selectedTeacher.name}</h2>
                    <div className="form-row">
                        <div>
                            <label>Добавить предмет</label>
                            <select value={assignSubjectId} onChange={e => setAssignSubjectId(Number(e.target.value))}>
                                <option value={0}>-- выберите --</option>
                                {availableSubjects.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                            </select>
                        </div>
                        <button className="btn-primary" disabled={!assignSubjectId} onClick={handleAssign}>Назначить</button>
                        <button className="btn-secondary" onClick={() => setSelectedTeacher(null)}>Закрыть</button>
                    </div>
                    <table>
                        <thead><tr><th>Предмет</th><th>Действие</th></tr></thead>
                        <tbody>
                            {teacherSubjects.map(s => (
                                <tr key={s.id}>
                                    <td>{s.name}</td>
                                    <td><button className="btn-danger" onClick={() => handleRemoveSubject(s.id)}>Убрать</button></td>
                                </tr>
                            ))}
                            {teacherSubjects.length === 0 && <tr><td colSpan={2} style={{ color: '#9ca3af', textAlign: 'center' }}>Предметы не назначены</td></tr>}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
}
