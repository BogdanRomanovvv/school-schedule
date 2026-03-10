export interface Class {
    id: number;
    name: string;
}

export interface Subject {
    id: number;
    name: string;
}

export interface Teacher {
    id: number;
    name: string;
    max_hours_per_week: number;
}

export interface Curriculum {
    id: number;
    class_id: number;
    subject_id: number;
    hours_per_week: number;
}

export interface ScheduleEntry {
    id: number;
    class_id: number;
    subject_id: number;
    teacher_id: number;
    day: number;
    lesson_number: number;
    class_name: string;
    subject_name: string;
    teacher_name: string;
}

export const DAY_NAMES = ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница'];
export const LESSON_NUMBERS = [1, 2, 3, 4, 5, 6, 7];
