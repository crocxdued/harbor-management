import { useState, useEffect } from 'react'
import { getUsers, createUser, deleteUser } from '../api/api'
import { useAuth } from '../context/AuthContext'
import Modal from '../components/Modal'
import s from '../styles/page.module.css'

const ROLES = ['admin','dispatcher','operator']
const R_LBL  = { admin:'Администратор', dispatcher:'Диспетчер', operator:'Оператор' }
const R_BG   = { admin:'#fee2e2', dispatcher:'#eff6ff', operator:'#f0fdf4' }
const R_CLR  = { admin:'#dc2626', dispatcher:'#1d4ed8', operator:'#16a34a' }
const EMPTY  = { name:'', email:'', age:'', role:'operator', password:'' }

export default function Users() {
  const { user: me, hasRole } = useAuth()
  const [list,   setList]   = useState([])
  const [busy,   setBusy]   = useState(true)
  const [err,    setErr]    = useState('')
  const [modal,  setModal]  = useState(false)
  const [form,   setForm]   = useState(EMPTY)
  const [saving, setSaving] = useState(false)

  if (!hasRole('admin')) return <div style={{padding:'2rem',color:'#dc2626'}}>Доступ только для администраторов.</div>

  const load = () => { setBusy(true); getUsers().then(setList).catch(e=>setErr(e.message)).finally(()=>setBusy(false)) }
  useEffect(load, [])

  const create = async e => {
    e.preventDefault(); setSaving(true)
    try {
      await createUser({ ...form, age: form.age ? parseInt(form.age) : undefined })
      setModal(false); load()
    } catch(e) { setErr(e.message) } finally { setSaving(false) }
  }

  const del = async id => {
    if (id === me?.id) { setErr('Нельзя удалить собственный аккаунт'); return }
    if (!confirm('Удалить пользователя?')) return
    try { await deleteUser(id); load() } catch(e) { setErr(e.message) }
  }

  return (
    <div>
      <div className={s.hdr}>
        <div><h1 className={s.title}>Пользователи</h1><p className={s.sub}>{list.length} в системе</p></div>
        <button className={s.btnPri} onClick={()=>{ setForm(EMPTY); setErr(''); setModal(true) }}>+ Добавить</button>
      </div>

      {err && <div className={s.err}>{err}</div>}
      {busy ? <div className={s.loading}>Загрузка...</div> : (
        <div className={s.card}>
          <table className={s.tbl}>
            <thead><tr><th>#</th><th>Имя</th><th>Email</th><th>Возраст</th><th>Роль</th><th>Дата</th><th></th></tr></thead>
            <tbody>
              {list.map(u => (
                <tr key={u.id}>
                  <td className={s.mono}>{u.id}</td>
                  <td><strong>{u.name}</strong>{u.id===me?.id&&<span style={{fontSize:'.7rem',color:'#9ca3af',marginLeft:'.35rem'}}>(вы)</span>}</td>
                  <td>{u.email}</td>
                  <td>{u.age ?? '—'}</td>
                  <td><span style={{display:'inline-block',padding:'3px 9px',borderRadius:'5px',fontSize:'.74rem',fontWeight:600,background:R_BG[u.role],color:R_CLR[u.role]}}>{R_LBL[u.role]}</span></td>
                  <td>{new Date(u.created_at).toLocaleDateString('ru-RU')}</td>
                  <td>{u.id!==me?.id&&<button className={s.btnDel} onClick={()=>del(u.id)}>Удалить</button>}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {modal && (
        <Modal title="Добавить пользователя" onClose={()=>setModal(false)}>
          <form onSubmit={create} className={s.form}>
            <div className={s.row}>
              <div className={s.fld}><label>Имя *</label><input required value={form.name} onChange={e=>setForm({...form,name:e.target.value})} placeholder="Иван Петров"/></div>
              <div className={s.fld}><label>Возраст</label><input type="number" min="18" max="120" value={form.age} onChange={e=>setForm({...form,age:e.target.value})} placeholder="30"/></div>
            </div>
            <div className={s.fld}><label>Email *</label><input type="email" required value={form.email} onChange={e=>setForm({...form,email:e.target.value})} placeholder="user@harbor.ru"/></div>
            <div className={s.row}>
              <div className={s.fld}><label>Роль *</label>
                <select value={form.role} onChange={e=>setForm({...form,role:e.target.value})}>
                  {ROLES.map(r=><option key={r} value={r}>{R_LBL[r]}</option>)}
                </select>
              </div>
              <div className={s.fld}><label>Пароль *</label><input type="password" required minLength={6} value={form.password} onChange={e=>setForm({...form,password:e.target.value})} placeholder="Мин. 6 символов"/></div>
            </div>
            {err && <div className={s.ferr}>{err}</div>}
            <div className={s.fbtns}>
              <button type="button" className={s.btnSec} onClick={()=>setModal(false)}>Отмена</button>
              <button type="submit" className={s.btnPri} disabled={saving}>{saving?'Создание...':'Создать'}</button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  )
}
