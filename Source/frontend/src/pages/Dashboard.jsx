import { useState, useEffect } from 'react'
import { getShips, getVisits } from '../api/api'
import { useAuth } from '../context/AuthContext'
import s from '../styles/page.module.css'
import d from './Dashboard.module.css'

const STATUS = { planned:'Запланирован', active:'Активен', completed:'Завершён', cancelled:'Отменён' }
const ST_CLS  = { planned: s.sPlanned, active: s.sActive, completed: s.sDone, cancelled: s.sCancelled }

export default function Dashboard() {
  const { user } = useAuth()
  const [ships,  setShips]  = useState([])
  const [visits, setVisits] = useState([])
  const [busy,   setBusy]   = useState(true)

  useEffect(() => {
    Promise.all([getShips(), getVisits()])
      .then(([sh, vi]) => { setShips(sh); setVisits(vi) })
      .finally(() => setBusy(false))
  }, [])

  const active = visits.filter(v => ['active','planned'].includes(v.status))
  const recent = [...visits].slice(0, 6)

  return (
    <div>
      <div className={s.hdr}>
        <div>
          <h1 className={s.title}>Дашборд</h1>
          <p className={s.sub}>Добро пожаловать, {user?.name}</p>
        </div>
      </div>

      <div className={d.grid}>
        {[
          { icon:'🚢', val: ships.length,  lbl:'Судов в реестре' },
          { icon:'📋', val: visits.length, lbl:'Всего визитов' },
          { icon:'⚡', val: active.length, lbl:'Активных визитов', hi: true },
        ].map(c => (
          <div key={c.lbl} className={`${d.stat} ${c.hi ? d.hi : ''}`}>
            <span className={d.icon}>{c.icon}</span>
            <span className={d.val}>{c.val}</span>
            <span className={d.lbl}>{c.lbl}</span>
          </div>
        ))}
      </div>

      <div className={s.card} style={{marginTop:'1.5rem'}}>
        <div style={{padding:'.85rem 1rem',borderBottom:'1px solid #f3f4f6',fontWeight:600,fontSize:'.9rem',color:'#111'}}>
          Последние визиты
        </div>
        {busy ? <div className={s.loading}>Загрузка...</div> : (
          <table className={s.tbl}>
            <thead><tr><th>Судно</th><th>Причал</th><th>Прибытие</th><th>Цель</th><th>Статус</th></tr></thead>
            <tbody>
              {recent.length === 0
                ? <tr><td colSpan={5} className={s.empty}>Нет данных</td></tr>
                : recent.map(v => (
                  <tr key={v.id}>
                    <td>{v.ship_name  || `#${v.ship_id}`}</td>
                    <td>{v.berth_number || `#${v.berth_id}`}</td>
                    <td>{new Date(v.arrival_time).toLocaleString('ru-RU',{day:'2-digit',month:'2-digit',year:'numeric',hour:'2-digit',minute:'2-digit'})}</td>
                    <td>{v.purpose}</td>
                    <td><span className={`${s.pill} ${ST_CLS[v.status]}`}>{STATUS[v.status]}</span></td>
                  </tr>
                ))
              }
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
