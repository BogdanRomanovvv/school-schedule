package domain

// Class — школьный класс
type Class struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

// Subject — учебный предмет
type Subject struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

// Teacher — учитель
type Teacher struct {
	ID              int    `db:"id"              json:"id"`
	Name            string `db:"name"            json:"name"`
	MaxHoursPerWeek int    `db:"max_hours_per_week" json:"max_hours_per_week"`
	// HomeroomClassID != nil означает классного руководителя — учитель ведёт уроки ТОЛЬКО в этом классе
	HomeroomClassID *int `db:"homeroom_class_id" json:"homeroom_class_id,omitempty"`
}

// TeacherSubject — связь учитель ↔ предмет
type TeacherSubject struct {
	TeacherID int `db:"teacher_id" json:"teacher_id"`
	SubjectID int `db:"subject_id" json:"subject_id"`
}

// Curriculum — учебный план: сколько часов в неделю предмет идёт у класса
type Curriculum struct {
	ID           int `db:"id"            json:"id"`
	ClassID      int `db:"class_id"      json:"class_id"`
	SubjectID    int `db:"subject_id"    json:"subject_id"`
	HoursPerWeek int `db:"hours_per_week" json:"hours_per_week"`
}

// ScheduleEntry — одна запись итогового расписания
type ScheduleEntry struct {
	ID           int `db:"id"           json:"id"`
	ClassID      int `db:"class_id"     json:"class_id"`
	SubjectID    int `db:"subject_id"   json:"subject_id"`
	TeacherID    int `db:"teacher_id"   json:"teacher_id"`
	Day          int `db:"day"          json:"day"`            // 0=ПН … 4=ПТ
	LessonNumber int `db:"lesson_number" json:"lesson_number"` // 0–6
}

// ScheduleEntryRich — расписание с человекочитаемыми именами
type ScheduleEntryRich struct {
	ScheduleEntry
	ClassName   string `db:"class_name"   json:"class_name"`
	SubjectName string `db:"subject_name" json:"subject_name"`
	TeacherName string `db:"teacher_name" json:"teacher_name"`
}
