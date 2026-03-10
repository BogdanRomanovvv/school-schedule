import { useEffect, useState } from 'react';
import * as api from '../api';
import type { Class } from '../types';

export default function ClassesPage() {
    const [classes, setClasses] = useState<Class[]>([]);
    const [name, setName] = useState('');
    const [editing, setEditing] = useState<Class | null>(null);
    const [error, setError] = useState('');

    const load = async () => {
        try { setClasses(await api.getClasses()); }
        catch (e: unknown) { setError(String(e)); }
    };

    useEffect(() => { load(); }, []);

    const handleSave = async () => {
        if (!name.trim()) return;
        try {
            if (editing) {
                await api.updateClass(editing.id, { name });
            } else {
                await api.createClass({ name });
            }
            setName(''); setEditing(null);
            load();
        } catch (e: unknown) { setError(String(e)); }
    };

    const handleEdit = (c: Class) => { setEditing(c); setName(c.name); };

    const handleDelete = async (id: number) => {
        if (!confirm('Удалить класс?')) return;
        try { await api.deleteClass(id); load(); }
        catch (e: unknown) { setError(String(e)); }
    };

    return (
        <div>
            <div className="card">
                <h2>Классы</h2>
                {error && <div className="alert alert-error">{error}</div>}
                <div className="form-row">
                    <div>
                        <label>Название класса</label>
                        <input
                            value={name}
                            onChange={e => setName(e.target.value)}
                            placeholder="например: 5А"
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
                        {classes.map((c, idx) => (
                            <tr key={c.id}>
                                <td>{idx + 1}</td>
                                <td>{c.name}</td>
                                <td>
                                    <div className="action-btns">
                                        <button className="btn-secondary" onClick={() => handleEdit(c)}>Изменить</button>
                                        <button className="btn-danger" onClick={() => handleDelete(c.id)}>Удалить</button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                        {classes.length === 0 && (
                            <tr><td colSpan={3} style={{ textAlign: 'center', color: '#9ca3af' }}>Нет данных</td></tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
