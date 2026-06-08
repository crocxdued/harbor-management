import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import s from './Layout.module.css'

const ROLE_LABEL = { admin: 'Администратор', dispatcher: 'Диспетчер', operator: 'Оператор' }
const ROLE_COLOR = { admin: '#ef4444', dispatcher: '#3b82f6', operator: '#22c55e' }

export default function Layout() {
  const { user, logout, hasRole } = useAuth()
  const nav = useNavigate()

  return (
    <div className={s.wrap}>
      <aside className={s.side}>
        <div className={s.brand}>
          <span>⚓</span>
          <div>
            <div className={s.brandName}>Harbor MS</div>
            <div className={s.brandSub}>Управление портом</div>
          </div>
        </div>

        <nav className={s.nav}>
          {[
            { to: '/',        label: 'Дашборд', icon: '📊', exact: true },
            { to: '/ships',   label: 'Суда',    icon: '🚢' },
            { to: '/visits',  label: 'Визиты',  icon: '📋' },
            ...(hasRole('admin') ? [{ to: '/users', label: 'Пользователи', icon: '👤' }] : []),
          ].map(({ to, label, icon, exact }) => (
            <NavLink key={to} to={to} end={exact}
              className={({ isActive }) => `${s.link} ${isActive ? s.active : ''}`}>
              <span>{icon}</span>{label}
            </NavLink>
          ))}
        </nav>

        <div className={s.foot}>
          <div className={s.uName}>{user?.name}</div>
          <div className={s.uRole} style={{ color: ROLE_COLOR[user?.role] }}>{ROLE_LABEL[user?.role]}</div>
          <div className={s.uEmail}>{user?.email}</div>
          <button className={s.logout} onClick={() => { logout(); nav('/login') }}>Выйти</button>
        </div>
      </aside>

      <main className={s.main}>
        <Outlet />
      </main>
    </div>
  )
}
