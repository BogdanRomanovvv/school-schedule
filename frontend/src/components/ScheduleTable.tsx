import { useState } from 'react';
import type { ScheduleEntry } from '../types';

interface Props {
    entries: ScheduleEntry[];
    dayNames: string[];
    lessonNumbers: number[];
    onEntryMove: (entry: ScheduleEntry) => void;
}

export default function ScheduleTable({ entries, dayNames, lessonNumbers, onEntryMove }: Props) {
    const [dragEntry, setDragEntry] = useState<ScheduleEntry | null>(null);
    const [editEntry, setEditEntry] = useState<ScheduleEntry | null>(null);

    // Build grid: day → lessonIndex → entry
    const grid: Record<number, Record<number, ScheduleEntry>> = {};
    for (const e of entries) {
        if (!grid[e.day]) grid[e.day] = {};
        grid[e.day][e.lesson_number] = e;
    }

    const handleDragStart = (e: React.DragEvent, entry: ScheduleEntry) => {
        setDragEntry(entry);
        e.dataTransfer.effectAllowed = 'move';
    };

    const handleDrop = (e: React.DragEvent, day: number, lessonIndex: number) => {
        e.preventDefault();
        if (!dragEntry) return;
        if (dragEntry.day === day && dragEntry.lesson_number === lessonIndex) return;

        // Check if target slot is occupied
        if (grid[day]?.[lessonIndex]) {
            alert('Этот слот уже занят другим уроком');
            return;
        }

        onEntryMove({ ...dragEntry, day, lesson_number: lessonIndex });
        setDragEntry(null);
    };

    const handleDragOver = (e: React.DragEvent) => { e.preventDefault(); e.dataTransfer.dropEffect = 'move'; };

    return (
        <div className="schedule-grid">
            <table>
                <thead>
                    <tr>
                        <th style={{ width: 60 }}>Урок</th>
                        {dayNames.map(d => <th key={d}>{d}</th>)}
                    </tr>
                </thead>
                <tbody>
                    {lessonNumbers.map((lessonNum, lessonIdx) => (
                        <tr key={lessonNum}>
                            <td style={{ fontWeight: 600, textAlign: 'center', color: '#6b7280' }}>{lessonNum}</td>
                            {dayNames.map((_, dayIdx) => {
                                const entry = grid[dayIdx]?.[lessonIdx];
                                return (
                                    <td
                                        key={dayIdx}
                                        className="schedule-cell"
                                        onDragOver={handleDragOver}
                                        onDrop={e => handleDrop(e, dayIdx, lessonIdx)}
                                        style={{ background: dragEntry ? '#f0f9ff' : undefined }}
                                    >
                                        {entry ? (
                                            <div
                                                className="lesson-card"
                                                draggable
                                                onDragStart={e => handleDragStart(e, entry)}
                                                onClick={() => setEditEntry(entry)}
                                                title="Кликните для просмотра / перетащите для перемещения"
                                            >
                                                <strong>{entry.subject_name}</strong>
                                                <span style={{ color: '#4b5563', fontSize: '.78rem' }}>
                                                    {entry.class_name && entry.teacher_name
                                                        ? `${entry.class_name} · ${entry.teacher_name}`
                                                        : entry.class_name || entry.teacher_name}
                                                </span>
                                            </div>
                                        ) : null}
                                    </td>
                                );
                            })}
                        </tr>
                    ))}
                </tbody>
            </table>

            {/* Detail modal */}
            {editEntry && (
                <div style={{
                    position: 'fixed', inset: 0, background: 'rgba(0,0,0,.4)',
                    display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 100,
                }} onClick={() => setEditEntry(null)}>
                    <div style={{
                        background: '#fff', borderRadius: 12, padding: 28, minWidth: 300, boxShadow: '0 8px 32px rgba(0,0,0,.2)',
                    }} onClick={e => e.stopPropagation()}>
                        <h3 style={{ marginBottom: 16 }}>Детали урока</h3>
                        <table style={{ width: '100%' }}>
                            <tbody>
                                <tr><td style={{ fontWeight: 600, paddingRight: 12 }}>Класс</td><td>{editEntry.class_name}</td></tr>
                                <tr><td style={{ fontWeight: 600 }}>Предмет</td><td>{editEntry.subject_name}</td></tr>
                                <tr><td style={{ fontWeight: 600 }}>Учитель</td><td>{editEntry.teacher_name}</td></tr>
                                <tr><td style={{ fontWeight: 600 }}>День</td><td>{'ПНВТСРЧТПТ'.slice(editEntry.day * 2, editEntry.day * 2 + 2)}</td></tr>
                                <tr><td style={{ fontWeight: 600 }}>Урок №</td><td>{editEntry.lesson_number + 1}</td></tr>
                            </tbody>
                        </table>
                        <button className="btn-secondary" style={{ marginTop: 20 }} onClick={() => setEditEntry(null)}>Закрыть</button>
                    </div>
                </div>
            )}
        </div>
    );
}
