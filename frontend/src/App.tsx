import { BrowserRouter, NavLink, Route, Routes, Navigate } from 'react-router-dom';
import ClassesPage from './pages/ClassesPage';
import SubjectsPage from './pages/SubjectsPage';
import TeachersPage from './pages/TeachersPage';
import CurriculumPage from './pages/CurriculumPage';
import SchedulePage from './pages/SchedulePage';

const nav = [
    { to: '/classes', label: '📚 Классы' },
    { to: '/subjects', label: '📖 Предметы' },
    { to: '/teachers', label: '👩‍🏫 Учителя' },
    { to: '/curriculum', label: '📋 Учебный план' },
    { to: '/schedule', label: '🗓 Расписание' },
];

export default function App() {
    return (
        <BrowserRouter>
            <div className="layout">
                <aside className="sidebar">
                    <h1>📅 Расписание</h1>
                    <nav>
                        {nav.map(n => (
                            <NavLink
                                key={n.to}
                                to={n.to}
                                className={({ isActive }) => isActive ? 'active' : ''}
                            >
                                {n.label}
                            </NavLink>
                        ))}
                    </nav>
                </aside>
                <main className="main">
                    <Routes>
                        <Route path="/" element={<Navigate to="/schedule" replace />} />
                        <Route path="/classes" element={<ClassesPage />} />
                        <Route path="/subjects" element={<SubjectsPage />} />
                        <Route path="/teachers" element={<TeachersPage />} />
                        <Route path="/curriculum" element={<CurriculumPage />} />
                        <Route path="/schedule" element={<SchedulePage />} />
                    </Routes>
                </main>
            </div>
        </BrowserRouter>
    );
}
