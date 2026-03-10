import { useEffect, useState } from 'react';
import * as api from '../api';
import type { Class, Teacher, ScheduleEntry } from '../types';
import { DAY_NAMES, LESSON_NUMBERS } from '../types';
import ScheduleTable from '../components/ScheduleTable';

type ViewMode = 'class' | 'teacher' | 'day';

export default function SchedulePage() {
    const [classes, setClasses] = useState<Class[]>([]);
    const [teachers, setTeachers] = useState<Teacher[]>([]);
    const [entries, setEntries] = useState<ScheduleEntry[]>([]);
    const [allEntries, setAllEntries] = useState<ScheduleEntry[]>([]);
    const [viewMode, setViewMode] = useState<ViewMode>('class');
    const [selectedId, setSelectedId] = useState<number>(0);
    const [selectedDay, setSelectedDay] = useState<number>(0);
    const [loading, setLoading] = useState(false);
    const [status, setStatus] = useState<{ type: 'success' | 'error' | 'info'; msg: string } | null>(null);

    useEffect(() => {
        Promise.all([api.getClasses(), api.getTeachers()]).then(([cls, tch]) => {
            setClasses(cls); setTeachers(tch);
            if (cls.length) setSelectedId(cls[0].id);
        });
        api.getSchedule().then(setAllEntries);
    }, []);

    useEffect(() => {
        if (viewMode === 'day') {
            api.getSchedule().then(setAllEntries);
            return;
        }
        if (!selectedId) return;
        loadEntries();
    }, [selectedId, viewMode]);

    const loadEntries = async () => {
        try {
            const data = viewMode === 'class'
                ? await api.getScheduleByClass(selectedId)
                : await api.getScheduleByTeacher(selectedId);
            setEntries(data);
        } catch {
            setEntries([]);
        }
    };

    const handleGenerate = async () => {
        setLoading(true);
        setStatus({ type: 'info', msg: 'Генерация расписания...' });
        try {
            const result = await api.generateSchedule();
            setStatus({ type: 'success', msg: `${result.message} (${result.count} уроков)` });
            const all = await api.getSchedule();
            setAllEntries(all);
            if (viewMode !== 'day') loadEntries();
        } catch (e: unknown) {
            const msg = (e as { response?: { data?: { error?: string } } })?.response?.data?.error ?? String(e);
            setStatus({ type: 'error', msg });
        } finally {
            setLoading(false);
        }
    };

    const handleClear = async () => {
        if (!confirm('Удалить всё расписание?')) return;
        await api.clearSchedule();
        setEntries([]);
        setAllEntries([]);
        setStatus({ type: 'info', msg: 'Расписание очищено' });
    };

    const handleEntryMove = async (entry: ScheduleEntry) => {
        try {
            await api.updateScheduleEntry(entry);
            const all = await api.getSchedule();
            setAllEntries(all);
            if (viewMode !== 'day') loadEntries();
        } catch (e: unknown) {
            const msg = (e as { response?: { data?: { error?: string } } })?.response?.data?.error ?? String(e);
            setStatus({ type: 'error', msg });
        }
    };

    // ─── Вид "По дням" ──────────────────────────────────────────────────────────
    const dayEntries = allEntries.filter(e => e.day === selectedDay);
    // Уникальные классы, у которых есть уроки в этот день (сортируем по имени)
    const dayClassIds = [...new Set(dayEntries.map(e => e.class_id ?? 0))].filter(Boolean);
    const dayClasses = dayClassIds
        .map(id => classes.find(c => c.id === id))
        .filter(Boolean) as Class[];
    dayClasses.sort((a, b) => a.name.localeCompare(b.name, 'ru', { numeric: true }));

    const maxLesson = dayEntries.reduce((m, e) => Math.max(m, e.lesson_number), -1);
    const lessonRows = maxLesson >= 0 ? Array.from({ length: maxLesson + 1 }, (_, i) => i) : [];

    const dayGrid: Record<number, Record<number, ScheduleEntry>> = {};
    for (const e of dayEntries) {
        if (!dayGrid[e.lesson_number]) dayGrid[e.lesson_number] = {};
        if (e.class_id) dayGrid[e.lesson_number][e.class_id] = e;
    }

    return (
        <div>
            {/* Controls */}
            <div className="card">
                <h2>Расписание</h2>

                <div className="form-row" style={{ marginBottom: 12 }}>
                    <button className="btn-primary" onClick={handleGenerate} disabled={loading}>
                        {loading ? <><span className="spinner" style={{ marginRight: 8 }} />Генерация...</> : '⚡ Сгенерировать расписание'}
                    </button>
                    <button className="btn-danger" onClick={handleClear}>🗑 Очистить</button>
                    <button className="btn-success" onClick={() => api.exportByClass()}>📥 Excel по классам</button>
                    <button className="btn-success" onClick={() => api.exportByTeacher()}>📥 Excel по учителям</button>
                </div>

                {status && (
                    <div className={`alert alert-${status.type}`}>{status.msg}</div>
                )}

                {/* View selector */}
                <div className="tabs">
                    <div className={`tab ${viewMode === 'class' ? 'active' : ''}`} onClick={() => { setViewMode('class'); if (classes.length) setSelectedId(classes[0].id); }}>
                        По классу
                    </div>
                    <div className={`tab ${viewMode === 'teacher' ? 'active' : ''}`} onClick={() => { setViewMode('teacher'); if (teachers.length) setSelectedId(teachers[0].id); }}>
                        По учителю
                    </div>
                    <div className={`tab ${viewMode === 'day' ? 'active' : ''}`} onClick={() => setViewMode('day')}>
                        По дням (все классы)
                    </div>
                </div>

                {viewMode !== 'day' && (
                    <div className="form-row">
                        <div>
                            <label>{viewMode === 'class' ? 'Класс' : 'Учитель'}</label>
                            <select value={selectedId} onChange={e => setSelectedId(Number(e.target.value))}>
                                {(viewMode === 'class' ? classes : teachers).map(item => (
                                    <option key={item.id} value={item.id}>{item.name}</option>
                                ))}
                            </select>
                        </div>
                    </div>
                )}

                {viewMode === 'day' && (
                    <div className="form-row">
                        <div>
                            <label>День недели</label>
                            <select value={selectedDay} onChange={e => setSelectedDay(Number(e.target.value))}>
                                {DAY_NAMES.map((d, i) => <option key={i} value={i}>{d}</option>)}
                            </select>
                        </div>
                    </div>
                )}
            </div>

            {/* Grid */}
            {viewMode === 'day' ? (
                <div className="card">
                    <h3 style={{ marginBottom: 16 }}>{DAY_NAMES[selectedDay]}</h3>
                    {dayClasses.length === 0 ? (
                        <div className="alert alert-info">Нет уроков в этот день.</div>
                    ) : (
                        <div className="schedule-grid" style={{ overflowX: 'auto' }}>
                            <table>
                                <thead>
                                    <tr>
                                        <th style={{ width: 60 }}>Урок</th>
                                        {dayClasses.map(c => <th key={c.id}>{c.name}</th>)}
                                    </tr>
                                </thead>
                                <tbody>
                                    {lessonRows.map(lessonIdx => (
                                        <tr key={lessonIdx}>
                                            <td style={{ fontWeight: 600, textAlign: 'center', color: '#6b7280' }}>
                                                {LESSON_NUMBERS[lessonIdx]}
                                            </td>
                                            {dayClasses.map(c => {
                                                const e = dayGrid[lessonIdx]?.[c.id];
                                                return (
                                                    <td key={c.id} className="schedule-cell">
                                                        {e ? (
                                                            <div className="lesson-card">
                                                                <strong>{e.subject_name}</strong>
                                                                <span style={{ color: '#4b5563', fontSize: '.78rem' }}>{e.teacher_name}</span>
                                                            </div>
                                                        ) : null}
                                                    </td>
                                                );
                                            })}
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </div>
            ) : entries.length > 0 ? (
                <div className="card">
                    <ScheduleTable
                        entries={entries}
                        dayNames={DAY_NAMES}
                        lessonNumbers={LESSON_NUMBERS}
                        onEntryMove={handleEntryMove}
                    />
                </div>
            ) : (
                <div className="card">
                    <div className="alert alert-info">Расписание не сгенерировано или нет данных для выбранного фильтра.</div>
                </div>
            )}
        </div>
    );
}
