import axios from 'axios';
import type { Class, Subject, Teacher, Curriculum, ScheduleEntry } from '../types';

const api = axios.create({ baseURL: '/api' });

// ─── Classes ──────────────────────────────────────────────────────────────────
export const getClasses = () => api.get<Class[]>('/classes').then(r => r.data ?? []);
export const createClass = (d: Omit<Class, 'id'>) => api.post<Class>('/classes', d).then(r => r.data);
export const updateClass = (id: number, d: Omit<Class, 'id'>) => api.put<Class>(`/classes/${id}`, d).then(r => r.data);
export const deleteClass = (id: number) => api.delete(`/classes/${id}`);

// ─── Subjects ────────────────────────────────────────────────────────────────
export const getSubjects = () => api.get<Subject[]>('/subjects').then(r => r.data ?? []);
export const createSubject = (d: Omit<Subject, 'id'>) => api.post<Subject>('/subjects', d).then(r => r.data);
export const updateSubject = (id: number, d: Omit<Subject, 'id'>) => api.put<Subject>(`/subjects/${id}`, d).then(r => r.data);
export const deleteSubject = (id: number) => api.delete(`/subjects/${id}`);

// ─── Teachers ────────────────────────────────────────────────────────────────
export const getTeachers = () => api.get<Teacher[]>('/teachers').then(r => r.data ?? []);
export const createTeacher = (d: Omit<Teacher, 'id'>) => api.post<Teacher>('/teachers', d).then(r => r.data);
export const updateTeacher = (id: number, d: Omit<Teacher, 'id'>) => api.put<Teacher>(`/teachers/${id}`, d).then(r => r.data);
export const deleteTeacher = (id: number) => api.delete(`/teachers/${id}`);
export const getTeacherSubjects = (teacherID: number) => api.get<Subject[]>(`/teachers/${teacherID}/subjects`).then(r => r.data ?? []);
export const assignSubject = (teacherID: number, subjectID: number) => api.post(`/teachers/${teacherID}/subjects/${subjectID}`);
export const removeSubject = (teacherID: number, subjectID: number) => api.delete(`/teachers/${teacherID}/subjects/${subjectID}`);

// ─── Curriculum ───────────────────────────────────────────────────────────────
export const getCurriculum = (classID?: number) =>
    api.get<Curriculum[]>('/curriculum', classID ? { params: { class_id: classID } } : undefined).then(r => r.data ?? []);
export const upsertCurriculum = (d: Omit<Curriculum, 'id'>) => api.post<Curriculum>('/curriculum', d).then(r => r.data);
export const deleteCurriculum = (classID: number, subjectID: number) => api.delete(`/curriculum/${classID}/${subjectID}`);

// ─── Schedule ────────────────────────────────────────────────────────────────
export const generateSchedule = () => api.post<{ message: string; count: number }>('/schedule/generate').then(r => r.data);
export const getSchedule = () => api.get<ScheduleEntry[]>('/schedule').then(r => r.data ?? []);
export const getScheduleByClass = (classID: number) => api.get<ScheduleEntry[]>('/schedule/by-class', { params: { class_id: classID } }).then(r => r.data ?? []);
export const getScheduleByTeacher = (teacherID: number) => api.get<ScheduleEntry[]>('/schedule/by-teacher', { params: { teacher_id: teacherID } }).then(r => r.data ?? []);
export const updateScheduleEntry = (e: ScheduleEntry) => api.put<ScheduleEntry>('/schedule/entry', e).then(r => r.data);
export const clearSchedule = () => api.delete('/schedule');

export const exportByClass = () => window.open('/api/schedule/export/class', '_blank');
export const exportByTeacher = () => window.open('/api/schedule/export/teacher', '_blank');
