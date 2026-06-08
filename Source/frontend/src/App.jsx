import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import Login    from './pages/Login'
import Layout   from './components/Layout'
import Dashboard from './pages/Dashboard'
import Ships    from './pages/Ships'
import Visits   from './pages/Visits'
import Users    from './pages/Users'

const Guard = ({ children }) => {
  const { user, loading } = useAuth()
  if (loading) return <div style={{padding:'2rem',color:'#6b7280'}}>Загрузка...</div>
  return user ? children : <Navigate to="/login" replace />
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<Guard><Layout /></Guard>}>
        <Route index       element={<Dashboard />} />
        <Route path="ships"  element={<Ships />} />
        <Route path="visits" element={<Visits />} />
        <Route path="users"  element={<Users />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppRoutes />
      </BrowserRouter>
    </AuthProvider>
  )
}
