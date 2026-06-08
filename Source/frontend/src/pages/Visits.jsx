import { useState, useEffect } from 'react'
import { getVisits, getShips, createVisit, updateVisit, deleteVisit } from '../api/api'
import { useAuth } from '../context/AuthContext'
import Modal from '../components/Modal'
import s from '../styles/page.module.css'

const ST_LBL = { planned:'Запланирован', active:'Активен', completed:'Завершён', cancelled:'Отменён' }
const ST_CLS = { planned:s.sPlanned, active:s.sActive, completed:s.sDone, cancelled:s.sCancelled }
const STATUSES = ['planned','active','completed','cancelled']

// Причалы захардкожены — они есть в БД, но не имеют отдельного CRUD в API
const BERTHS = [
  {id:1,label:'A-01 — грузовой'},
  {id:2,label:'A-02 — грузовой'},
  {id:3,label:'B-01 — нефтяной'},
  {id:4,label:'B-02 — нефтяной'},
  {id:5,label:'C-01 — пассажирский'},
  {id:6,label:'D-01 — сухой dok'},
]

export default function Visits() {
  const { hasRole } = useAuth()
  const [visits, setVisits] = useState([])
  const [ships,  setShips]  = useState([])
  const [busy,   setBusy]   = useState(true)
  const [err,    setErr]    = useState('')
  const [modal,  setModal]  = useState(null) // 'create' | 'status'
  const [sel,    setSel]    = useState(null)
  const [form,   setForm]   = useState({ ship_id:'', berth_id:'1', arrival_time:'', purpose:'' })
  const [stForm, setStForm] = useState({ status:'active', departure_time:'' })
  const [saving, setSaving] = useState(false)

  const load = () => {
    setBusy(true)
    Promise.all([getVisits(), getShips()])
      .then(([v,sh]) => { setVisits(v); setShips(sh) })
      .catch(e=>setErr(e.message))
      .finally(()=>setBusy(false))
  }
  useEffect(load, [])

  const openCreate = () => { setForm({ ship_id:'', berth_id:'1', arrival_time:'', purpose:'' }); setErr(''); setModal('create') }
  const openStatus = v => { setSel(v); setStForm({ status:v.status, departure_time:'' }); setErr(''); setModal('status') }

  const create = async e => {
    e.preventDefault(); setSaving(true)
    try {
      await createVisit({ ship_id:parseInt(form.ship_id), berth_id:parseInt(form.berth_id), arrival_time:new Date(form.arrival_time).toISOString(), purpose:form.purpose })
      setModal(null); load()
    } catch(e) { setErr(e.message) } finally { setSaving(false) }
  }

  const updateStatus = async e => {
    e.preventDefault(); setSaving(true)
    try {
      const d = { status: stForm.status }
      if (stForm.departure_time) d.departure_time = new Date(stForm.departure_time).toISOString()
      await updateVisit(sel.id, d)
      setModal(null); load()
    } catch(e) { setErr(e.message) } finally { setSaving(false) }
  }

  const del = async id => {
    if (!confirm('Удалить визит?')) return
    try { await deleteVisit(id); load() } catch(e) { setErr(e.message) }
  }

  const fmt = dt => new Date(dt).toLocaleString('ru-RU',{day:'2-digit',month:'2-digit',year:'numeric',hour:'2-digit',minute:'2-digit'})
  const canEdit = hasRole('admin','dispatcher')
  const canDel  = hasRole('admin')

  return (
    <div>
      <div className={s.hdr}>
        <div><h1 className={s.title}>Визиты судов</h1><p className={s.sub}>{visits.length} записей</p></div>
        {canEdit && <button className={s.btnPri} onClick={openCreate}>+ Зарегистрировать</button>}
      </div>

      {err && <div className={s.err}>{err}</div>}
      {busy ? <div className={s.loading}>Загрузка...</div> : (
        <div className={s.card}>
          <table className={s.tbl}>
            <thead><tr><th>#</th><th>Судно</th><th>Причал</th><th>Прибытие</th><th>Убытие</th><th>Цель</th><th>Статус</th>{(canEdit||canDel)&&<th></th>}</tr></thead>
            <tbody>
              {visits.length === 0
                ? <tr><td colSpan={8} className={s.empty}>Нет визитов</td></tr>
                : visits.map(v => (
                  <tr key={v.id}>
                    <td className={s.mono}>#{v.id}</td>
                    <td><strong>{v.ship_name || `#${v.ship_id}`}</strong></td>
                    <td>{v.berth_number || `#${v.berth_id}`}</td>
                    <td>{fmt(v.arrival_time)}</td>
                    <td>{v.departure_time ? fmt(v.departure_time) : '—'}</td>
                    <td>{v.purpose}</td>
                    <td><span className={`${s.pill} ${ST_CLS[v.status]}`}>{ST_LBL[v.status]}</span></td>
                    {(canEdit||canDel) && (
                      <td className={s.acts}>
                        {canEdit && !['completed','cancelled'].includes(v.status) && <button className={s.btnEdit} onClick={()=>openStatus(v)}>Статус</button>}
                        {canDel  && <button className={s.btnDel} onClick={()=>del(v.id)}>Удалить</button>}
                      </td>
                    )}
                  </tr>
                ))
              }
            </tbody>
          </table>
        </div>
      )}

      {modal === 'create' && (
        <Modal title="Зарегистрировать визит" onClose={()=>setModal(null)}>
          <form onSubmit={create} className={s.form}>
            <div className={s.fld}><label>Судно *</label>
              <select required value={form.ship_id} onChange={e=>setForm({...form,ship_id:e.target.value})}>
                <option value="">— выберите судно —</option>
                {ships.map(sh=><option key={sh.id} value={sh.id}>{sh.name} ({sh.imo_number})</option>)}
              </select>
            </div>
            <div className={s.fld}><label>Причал *</label>
              <select value={form.berth_id} onChange={e=>setForm({...form,berth_id:e.target.value})}>
                {BERTHS.map(b=><option key={b.id} value={b.id}>{b.label}</option>)}
              </select>
            </div>
            <div className={s.fld}><label>Дата и время прибытия *</label>
              <input type="datetime-local" required value={form.arrival_time} onChange={e=>setForm({...form,arrival_time:e.target.value})}/>
            </div>
            <div className={s.fld}><label>Цель визита *</label>
              <input required value={form.purpose} onChange={e=>setForm({...form,purpose:e.target.value})} placeholder="Разгрузка контейнеров"/>
            </div>
            {err && <div className={s.ferr}>{err}</div>}
            <div className={s.fbtns}>
              <button type="button" className={s.btnSec} onClick={()=>setModal(null)}>Отмена</button>
              <button type="submit" className={s.btnPri} disabled={saving}>{saving?'Сохранение...':'Зарегистрировать'}</button>
            </div>
          </form>
        </Modal>
      )}

      {modal === 'status' && sel && (
        <Modal title={`Визит #${sel.id} — изменить статус`} onClose={()=>setModal(null)}>
          <form onSubmit={updateStatus} className={s.form}>
            <div className={s.fld}><label>Новый статус</label>
              <select value={stForm.status} onChange={e=>setStForm({...stForm,status:e.target.value})}>
                {STATUSES.map(st=><option key={st} value={st}>{ST_LBL[st]}</option>)}
              </select>
            </div>
            {stForm.status === 'completed' && (
              <div className={s.fld}><label>Время убытия</label>
                <input type="datetime-local" value={stForm.departure_time} onChange={e=>setStForm({...stForm,departure_time:e.target.value})}/>
              </div>
            )}
            {err && <div className={s.ferr}>{err}</div>}
            <div className={s.fbtns}>
              <button type="button" className={s.btnSec} onClick={()=>setModal(null)}>Отмена</button>
              <button type="submit" className={s.btnPri} disabled={saving}>{saving?'Сохранение...':'Обновить'}</button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  )
}
