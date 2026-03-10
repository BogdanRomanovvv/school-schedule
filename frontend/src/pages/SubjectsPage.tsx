import { useEffect, useState } from 'react';
import * as api from '../api';
import type { Subject } from '../types';

export default function SubjectsPage() {
    const [subjects, setSubjects] = useState<Subject[]>([]);
    const [name, setName] = useState('');
    const [editing, setEditing] = useState<Subject | null>(null);
    const [error, setError] = useState('');

    const load = async () => {
        try { setSubjects(await api.getSubjects()); }
        catch (e: unknown) { setError(String(e)); }
    };

    useEffect(() => { load(); }, []);

    const handleSave = async () => {
        if (!name.trim()) return;
        try {
            if (editing) {
                await api.updateSubject(editing.id, { name });
            } else {
                await api.createSubject({ name });
            }
            setName(''); setEditing(null); load();
        } catch (e: unknown) { setError(String(e)); }
    };

    const handleDelete = async (id: number) => {
        if (!confirm('Удалить предмет?')) return;
        try { await api.deleteSubject(id); load(); }
        catch (e: unknown) { setError(String(e)); }
    };

    return (
        <div>
            <div className="card">
                <h2>Предметы</h2>
                {error && <div className="alert alert-error">{error}</div>}
                <div className="form-row">
                    <div>
                        <label>Название предмета</label>
                        <input
                            value={name}
                            onChange={e => setName(e.target.value)}
                            placeholder="например: Математика"
                            onKeyDown={e => e.key === 'Enter' && handleSave()}
                        />
                    </div>
                    <button className="btn-primary" onClick={handleSave}>
                        {editing ? 'Сохранить' : 'Добавить'}
                    </button>
                    {editing && (
                        <button className="btn-secondary" onClick={() => { setEditing(null); setName(''); }}>
                            Отмена
                        </button>
                    )}
                </div>

                <table>
                    <thead>
                        <tr><th>#</th><th>Название</th><th>Действия</th></tr>
                    </thead>
                    <tbody>
                        {subjects.map((s, idx) => (
                            <tr key={s.id}>
                                <td>{idx + 1}</td>
                                <td>{s.name}</td>
                                <td>
                                    <div className="action-btns">
                                        <button className="btn-secondary" onClick={() => { setEditing(s); setName(s.name); }}>Изменить</button>
                                        <button className="btn-danger" onClick={() => handleDelete(s.id)}>Удалить</button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                        {subjects.length === 0 && (
                            <tr><td colSpan={3} style={{ textAlign: 'center', color: '#9ca3af' }}>Нет данных</td></tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
