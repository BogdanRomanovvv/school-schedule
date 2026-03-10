package scheduler

import (
	"fmt"
	"sort"

	"school-schedule/internal/domain"
)

const (
	Days    = 5
	Lessons = 7
)

type Generator struct {
	curricula       []domain.Curriculum
	teachers        []domain.Teacher
	teacherSubjects map[int][]int
	subjectTeachers map[int][]domain.Teacher
}

func NewGenerator(
	curricula []domain.Curriculum,
	teachers []domain.Teacher,
	teacherSubjects map[int][]int,
) *Generator {
	g := &Generator{
		curricula:       curricula,
		teachers:        teachers,
		teacherSubjects: teacherSubjects,
		subjectTeachers: make(map[int][]domain.Teacher),
	}
	for _, t := range teachers {
		for _, sid := range teacherSubjects[t.ID] {
			g.subjectTeachers[sid] = append(g.subjectTeachers[sid], t)
		}
	}
	return g
}

func (g *Generator) Generate() ([]domain.ScheduleEntry, error) {

	// Phase 1: assign one teacher per (class, subject)
	cur := make([]domain.Curriculum, len(g.curricula))
	copy(cur, g.curricula)
	sort.SliceStable(cur, func(i, j int) bool {
		ti := len(g.subjectTeachers[cur[i].SubjectID])
		tj := len(g.subjectTeachers[cur[j].SubjectID])
		if ti != tj {
			return ti < tj
		}
		return cur[i].HoursPerWeek > cur[j].HoursPerWeek
	})

	classSubjectTeacher := make(map[[2]int]int)
	teacherAssigned := make(map[int]int)

	for _, c := range cur {
		candidates := g.subjectTeachers[c.SubjectID]
		if len(candidates) == 0 {
			return nil, fmt.Errorf("no teachers for subject id=%d", c.SubjectID)
		}
		// Если для класса есть классный руководитель, способный вести этот предмет —
		// назначаем только его (исключаем остальных кандидатов).
		// Иначе исключаем классных руководителей чужих классов.
		var homeroomForClass *domain.Teacher
		for i := range candidates {
			t := &candidates[i]
			if t.HomeroomClassID != nil && *t.HomeroomClassID == c.ClassID {
				homeroomForClass = t
				break
			}
		}
		var filtered []domain.Teacher
		if homeroomForClass != nil {
			// принудительно назначаем классного руководителя
			filtered = []domain.Teacher{*homeroomForClass}
		} else {
			// обычный класс: исключаем классных руководителей чужих классов
			for _, t := range candidates {
				if t.HomeroomClassID != nil && *t.HomeroomClassID != c.ClassID {
					continue
				}
				filtered = append(filtered, t)
			}
		}
		if len(filtered) == 0 {
			return nil, fmt.Errorf("no eligible teachers for class id=%d subject id=%d", c.ClassID, c.SubjectID)
		}
		best, bestLoad := -1, int(^uint(0)>>1)
		for _, t := range filtered {
			load := teacherAssigned[t.ID]
			if load+c.HoursPerWeek <= t.MaxHoursPerWeek && load < bestLoad {
				best, bestLoad = t.ID, load
			}
		}
		if best == -1 {
			for _, t := range filtered {
				if load := teacherAssigned[t.ID]; load < bestLoad {
					best, bestLoad = t.ID, load
				}
			}
		}
		classSubjectTeacher[[2]int{c.ClassID, c.SubjectID}] = best
		teacherAssigned[best] += c.HoursPerWeek
	}

	// Phase 2: place lessons
	classTotalHours := make(map[int]int)
	for _, c := range g.curricula {
		classTotalHours[c.ClassID] += c.HoursPerWeek
	}
	classMaxPerDay := make(map[int]int)
	for classID, total := range classTotalHours {
		classMaxPerDay[classID] = (total + Days - 1) / Days
	}

	teacherMaxHours := make(map[int]int)
	for _, t := range g.teachers {
		teacherMaxHours[t.ID] = t.MaxHoursPerWeek
	}
	teacherHoursPlaced := make(map[int]int)
	classSlot := make(map[[2]int]bool)
	teacherSlot := make(map[[2]int]bool)
	classDayCount := make(map[[2]int]int)

	type cdsKey struct{ classID, day, subjectID int }
	classDaySubject := make(map[cdsKey]bool)
	// classLastSlot[{classID,day}] = highest lesson slot placed so far on that class-day.
	// Initialized to -1 (nothing placed). Used for contiguous lesson placement to avoid gaps.
	classLastSlot := make(map[[2]int]int)

	result := make([]domain.ScheduleEntry, 0, 300)

	for idx, c := range cur {
		teacherID := classSubjectTeacher[[2]int{c.ClassID, c.SubjectID}]
		maxPerDay := classMaxPerDay[c.ClassID]
		placed := 0
		startDay := idx % Days

		// teacherDayLoad: number of slots already assigned to this teacher per day
		// (used to prefer days where the teacher has more free slots => contiguous placement)
		teacherDayLoad := func(tID, d int) int {
			count := 0
			for l := 0; l < Lessons; l++ {
				if teacherSlot[[2]int{tID, d*Lessons + l}] {
					count++
				}
			}
			return count
		}

		// Build a day order that interleaves the default startDay rotation with
		// preferring days where the teacher has fewer lessons already (fewer gaps).
		dayOrder := func(pass int) []int {
			days := make([]int, Days)
			for i := range days {
				days[i] = (startDay + i) % Days
			}
			// Sort: valid days (not skipped by classDaySubject, respects maxPerDay on pass=0)
			// first, then within those prefer days where teacher has lower load.
			sort.SliceStable(days, func(a, b int) bool {
				da, db := days[a], days[b]
				// On pass 0 push over-full days to the back
				if pass == 0 {
					aFull := classDayCount[[2]int{c.ClassID, da}] >= maxPerDay
					bFull := classDayCount[[2]int{c.ClassID, db}] >= maxPerDay
					if aFull != bFull {
						return !aFull
					}
				}
				// Prefer days where teacher already has lessons (contiguous) but fewer total
				// => lower teacher load = more room without gaps
				return teacherDayLoad(teacherID, da) < teacherDayLoad(teacherID, db)
			})
			return days
		}

		for pass := 0; pass < 2 && placed < c.HoursPerWeek; pass++ {
			for _, day := range dayOrder(pass) {
				if placed >= c.HoursPerWeek {
					break
				}
				if classDaySubject[cdsKey{c.ClassID, day, c.SubjectID}] {
					continue
				}
				if pass == 0 && classDayCount[[2]int{c.ClassID, day}] >= maxPerDay {
					continue
				}

				// Determine the preferred starting lesson slot: right after the last placed
				// lesson for this class on this day, to keep lessons contiguous.
				lastSlot, hasLesson := classLastSlot[[2]int{c.ClassID, day}]
				prefStart := 0
				if hasLesson {
					prefStart = lastSlot + 1
				}

				tryLesson := func(fromL, toL int) bool {
					for lesson := fromL; lesson < toL; lesson++ {
						slot := day*Lessons + lesson
						cKey := [2]int{c.ClassID, slot}
						tKey := [2]int{teacherID, slot}
						if classSlot[cKey] || teacherSlot[tKey] {
							continue
						}
						if teacherHoursPlaced[teacherID] >= teacherMaxHours[teacherID] {
							continue
						}
						classSlot[cKey] = true
						teacherSlot[tKey] = true
						classDayCount[[2]int{c.ClassID, day}]++
						classDaySubject[cdsKey{c.ClassID, day, c.SubjectID}] = true
						teacherHoursPlaced[teacherID]++
						lk := [2]int{c.ClassID, day}
						if cur, ok := classLastSlot[lk]; !ok || lesson > cur {
							classLastSlot[lk] = lesson
						}
						result = append(result, domain.ScheduleEntry{
							ClassID:      c.ClassID,
							SubjectID:    c.SubjectID,
							TeacherID:    teacherID,
							Day:          day,
							LessonNumber: lesson,
						})
						placed++
						return true
					}
					return false
				}

				// Try preferred contiguous range first, then fall back to earlier slots.
				if placed < c.HoursPerWeek {
					if !tryLesson(prefStart, Lessons) {
						tryLesson(0, prefStart)
					}
				}
			}
		}

		if placed < c.HoursPerWeek {
			return nil, fmt.Errorf(
				"failed: class id=%d subject id=%d (%d/%d placed)",
				c.ClassID, c.SubjectID, placed, c.HoursPerWeek,
			)
		}
	}

	// Phase 3: compact
	result = compact(result, teacherSlot)
	return result, nil
}

// compact eliminates per-class daily gaps via two strategies:
//
//	A) Within a (class,day): move any later lesson to the gap slot if the
//	   teacher is free there.
//	B) Inter-day: move the LAST lesson of another day (same class) to fill
//	   the gap, if the teacher is free at the gap slot and the subject is
//	   not already present on the gap day.
func compact(entries []domain.ScheduleEntry, teacherSlot map[[2]int]bool) []domain.ScheduleEntry {
	type cdKey struct{ classID, day int }

	groups := make(map[cdKey][]int)
	for i, e := range entries {
		k := cdKey{e.ClassID, e.Day}
		groups[k] = append(groups[k], i)
	}

	daySubjects := func(classID, day int) map[int]bool {
		used := make(map[int]bool)
		for _, idx := range groups[cdKey{classID, day}] {
			used[entries[idx].SubjectID] = true
		}
		return used
	}

	changed := true
	for iter := 0; changed && iter < 500; iter++ {
		changed = false

		// Sort keys for deterministic processing order
		keys := make([]cdKey, 0, len(groups))
		for k := range groups {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(a, b int) bool {
			if keys[a].classID != keys[b].classID {
				return keys[a].classID < keys[b].classID
			}
			return keys[a].day < keys[b].day
		})

		for _, key := range keys {
			indices := groups[key]
			sort.Slice(indices, func(a, b int) bool {
				return entries[indices[a]].LessonNumber < entries[indices[b]].LessonNumber
			})
			groups[key] = indices

			// Build occupied-slot set for correct gap detection
			occupied := make(map[int]bool, len(indices))
			maxSlot := 0
			for _, idx := range indices {
				ln := entries[idx].LessonNumber
				occupied[ln] = true
				if ln > maxSlot {
					maxSlot = ln
				}
			}

			for gapPos := 0; gapPos < maxSlot; gapPos++ {
				if occupied[gapPos] {
					continue
				}

				// Strategy A: fill gap at gapPos with any lesson whose LN > gapPos and teacher is free
				filled := false
				for _, idx := range indices {
					le := &entries[idx]
					if le.LessonNumber <= gapPos {
						continue // can only pull from later slots
					}
					tKeyOld := [2]int{le.TeacherID, le.Day*Lessons + le.LessonNumber}
					tKeyNew := [2]int{le.TeacherID, le.Day*Lessons + gapPos}
					if !teacherSlot[tKeyNew] {
						delete(teacherSlot, tKeyOld)
						teacherSlot[tKeyNew] = true
						occupied[le.LessonNumber] = false
						occupied[gapPos] = true
						le.LessonNumber = gapPos
						changed = true
						filled = true
						break
					}
				}
				if filled {
					break
				}

				// Strategy B: take the LAST lesson of another day (safe: no new gap on source day)
				usedOnGapDay := daySubjects(key.classID, key.day)

				for srcDay := 0; srcDay < Days && !filled; srcDay++ {
					if srcDay == key.day {
						continue
					}
					srcKey := cdKey{key.classID, srcDay}
					srcIdx := groups[srcKey]
					if len(srcIdx) == 0 {
						continue
					}
					sort.Slice(srcIdx, func(a, b int) bool {
						return entries[srcIdx[a]].LessonNumber < entries[srcIdx[b]].LessonNumber
					})
					groups[srcKey] = srcIdx

					last := len(srcIdx) - 1
					le := &entries[srcIdx[last]]

					if usedOnGapDay[le.SubjectID] {
						continue
					}
					tKeyOld := [2]int{le.TeacherID, srcDay*Lessons + le.LessonNumber}
					tKeyNew := [2]int{le.TeacherID, key.day*Lessons + gapPos}
					if teacherSlot[tKeyNew] {
						continue
					}

					delete(teacherSlot, tKeyOld)
					teacherSlot[tKeyNew] = true
					le.Day = key.day
					le.LessonNumber = gapPos

					groups[srcKey] = srcIdx[:last]
					groups[key] = append(groups[key], srcIdx[last])
					changed = true
					filled = true
				}

				// Strategy C: cross-class swap.
				// For each lesson of classA in this (day, LN>gapPos) whose teacher T is busy at gapPos
				// because some other class B has T at (day, gapPos):
				//   Swap: classA lesson moves to gapPos, classB lesson moves to LN.
				// Since both lessons share the same teacher T, teacherSlot doesn't change.
				// Valid when classB has no lesson at slot LN and no subject duplication occurs.
				for _, idxA := range indices {
					if filled {
						break
					}
					leA := &entries[idxA]
					if leA.LessonNumber <= gapPos {
						continue
					}
					lnA := leA.LessonNumber
					tID := leA.TeacherID

					// Teacher T must be busy at gapPos (otherwise Strategy A covered it)
					if !teacherSlot[[2]int{tID, key.day*Lessons + gapPos}] {
						continue
					}

					// Find class B that occupies (teacher T, day, gapPos)
					idxB := -1
					for j := range entries {
						e := &entries[j]
						if e.TeacherID == tID && e.Day == key.day && e.LessonNumber == gapPos {
							idxB = j
							break
						}
					}
					if idxB < 0 {
						continue
					}
					leB := &entries[idxB]
					classB := leB.ClassID

					// classB must not already have a lesson at slot lnA
					classBHasLnA := false
					for _, idx2 := range groups[cdKey{classB, key.day}] {
						if entries[idx2].LessonNumber == lnA {
							classBHasLnA = true
							break
						}
					}
					if classBHasLnA {
						continue
					}

					// No subject duplication after swap
					usedA := daySubjects(key.classID, key.day)
					usedB := daySubjects(classB, key.day)
					if usedA[leB.SubjectID] || usedB[leA.SubjectID] {
						continue
					}

					// Perform swap: just update LessonNumbers; teacherSlot unchanged
					leA.LessonNumber = gapPos
					leB.LessonNumber = lnA
					occupied[gapPos] = true
					// lnA is now occupied by classB lesson; occupied map is for classA only
					// so lnA is no longer occupied by classA:
					occupied[lnA] = false

					changed = true
					filled = true
				}

				if changed {
					break
				}
				// gap unfillable; continue to next position
			}
		}
	}

	return entries
}
