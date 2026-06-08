import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import s from './Login.module.css'

const HINTS = [
  { role: 'admin',      label: 'Администратор', email: 'admin@harbor.ru',      pass: 'admin123' },
  { role: 'dispatcher', label: 'Диспетчер',     email: 'dispatcher@harbor.ru', pass: 'disp123'  },
  { role: 'operator',   label: 'Оператор',      email: 'operator@harbor.ru',   pass: 'oper123'  },
]

export default function Login() {
  const { login } = useAuth()
  const nav = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [err, setErr] = useState('')
  const [busy, setBusy] = useState(false)

  const submit = async e => {
    e.preventDefault(); setErr(''); setBusy(true)
    try { await login(email, password); nav('/') }
    catch (e) { setErr(e.message) }
    finally { setBusy(false) }
  }

  return (
    <div className={s.page}>
      <div className={s.card}>
        <div className={s.logo}>
          <span className={s.anchor}>⚓</span>
          <h1 className={s.title}>Система управления портом</h1>
          <p className={s.sub}>Harbor Management System</p>
        </div>

        <form onSubmit={submit} className={s.form}>
          <div className={s.fld}>
            <label>Email</label>
            <input type="email" value={email} onChange={e=>setEmail(e.target.value)} placeholder="admin@harbor.ru" required />
          </div>
          <div className={s.fld}>
            <label>Пароль</label>
            <input type="password" value={password} onChange={e=>setPassword(e.target.value)} placeholder="••••••••" required />
          </div>
          {err && <div className={s.err}>{err}</div>}
          <button type="submit" className={s.btn} disabled={busy}>{busy ? 'Вход...' : 'Войти'}</button>
        </form>

        <div className={s.hints}>
          <p className={s.hintsLbl}>Тестовые аккаунты</p>
          {HINTS.map(h => (
            <div key={h.role} className={s.hint} onClick={() => { setEmail(h.email); setPassword(h.pass) }}>
              <span className={`${s.badge} ${s[h.role]}`}>{h.label}</span>
              <span>{h.email}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
