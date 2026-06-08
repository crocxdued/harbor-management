import { useState, useEffect } from 'react'
import { getShips, createShip, updateShip, deleteShip } from '../api/api'
import { useAuth } from '../context/AuthContext'
import Modal from '../components/Modal'
import s from '../styles/page.module.css'

const TYPES = ['cargo','tanker','passenger','warship']
const TYPE_LBL = { cargo:'Грузовое', tanker:'Танкер', passenger:'Пассажирское', warship:'Военное' }
const TYPE_ICO = { cargo:'📦', tanker:'🛢️', passenger:'🛳️', warship:'⚔️' }
const EMPTY = { name:'', imo_number:'', ship_type:'cargo', flag_country:'', gross_tonnage:'', year_built:'' }

export default function Ships() {
  const { hasRole } = useAuth()
  const [list,    setList]    = useState([])
  const [busy,    setBusy]    = useState(true)
  const [err,     setErr]     = useState('')
  const [modal,   setModal]   = useState(false)
  const [editing, setEditing] = useState(null)
  const [form,    setForm]    = useState(EMPTY)
  const [saving,  setSaving]  = useState(false)

  const load = () => { setBusy(true); getShips().then(setList).catch(e=>setErr(e.message)).finally(()=>setBusy(false)) }
  useEffect(load, [])

  const open = (ship = null) => {
    setEditing(ship)
    setForm(ship ? { name:ship.name, imo_number:ship.imo_number, ship_type:ship.ship_type, flag_country:ship.flag_country, gross_tonnage:ship.gross_tonnage, year_built:ship.year_built } : EMPTY)
    setErr(''); setModal(true)
  }

  const save = async e => {
    e.preventDefault(); setSaving(true)
    try {
      const d = { ...form, gross_tonnage: parseInt(form.gross_tonnage), year_built: parseInt(form.year_built) }
      editing ? await updateShip(editing.id, d) : await createShip(d)
      setModal(false); load()
    } catch(e) { setErr(e.message) } finally { setSaving(false) }
  }

  const del = async id => {
    if (!confirm('Удалить судно?')) return
    try { await deleteShip(id); load() } catch(e) { setErr(e.message) }
  }

  const canEdit = hasRole('admin','dispatcher')
  const canDel  = hasRole('admin')

  return (
    <div>
      <div className={s.hdr}>
        <div><h1 className={s.title}>Суда</h1><p className={s.sub}>{list.length} в реестре</p></div>
        {canEdit && <button className={s.btnPri} onClick={()=>open()}>+ Добавить</button>}
      </div>

      {err && <div className={s.err}>{err}</div>}
      {busy ? <div className={s.loading}>Загрузка...</div> : (
        <div className={s.card}>
          <table className={s.tbl}>
            <thead><tr><th>Название</th><th>ИМО</th><th>Тип</th><th>Флаг</th><th>Тоннаж</th><th>Год</th>{(canEdit||canDel)&&<th></th>}</tr></thead>
            <tbody>
              {list.length === 0
                ? <tr><td colSpan={7} className={s.empty}>Нет судов</td></tr>
                : list.map(ship => (
                  <tr key={ship.id}>
                    <td><strong>{ship.name}</strong></td>
                    <td className={s.mono}>{ship.imo_number}</td>
                    <td><span className={s.chip}>{TYPE_ICO[ship.ship_type]} {TYPE_LBL[ship.ship_type]}</span></td>
                    <td>{ship.flag_country}</td>
                    <td>{ship.gross_tonnage.toLocaleString('ru')} GRT</td>
                    <td>{ship.year_built}</td>
                    {(canEdit||canDel) && (
                      <td className={s.acts}>
                        {canEdit && <button className={s.btnEdit} onClick={()=>open(ship)}>Изменить</button>}
                        {canDel  && <button className={s.btnDel}  onClick={()=>del(ship.id)}>Удалить</button>}
                      </td>
                    )}
                  </tr>
                ))
              }
            </tbody>
          </table>
        </div>
      )}

      {modal && (
        <Modal title={editing ? 'Редактировать судно' : 'Добавить судно'} onClose={()=>setModal(false)}>
          <form onSubmit={save} className={s.form}>
            <div className={s.row}>
              <div className={s.fld}><label>Название *</label><input required value={form.name} onChange={e=>setForm({...form,name:e.target.value})} placeholder="Северный Ветер"/></div>
              <div className={s.fld}><label>Номер ИМО *</label><input required value={form.imo_number} onChange={e=>setForm({...form,imo_number:e.target.value})} placeholder="IMO1234567" disabled={!!editing}/></div>
            </div>
            <div className={s.row}>
              <div className={s.fld}><label>Тип *</label>
                <select value={form.ship_type} onChange={e=>setForm({...form,ship_type:e.target.value})}>
                  {TYPES.map(t=><option key={t} value={t}>{TYPE_LBL[t]}</option>)}
                </select>
              </div>
              <div className={s.fld}><label>Страна флага *</label><input required value={form.flag_country} onChange={e=>setForm({...form,flag_country:e.target.value})} placeholder="Россия"/></div>
            </div>
            <div className={s.row}>
              <div className={s.fld}><label>Тоннаж GRT *</label><input type="number" required min="1" value={form.gross_tonnage} onChange={e=>setForm({...form,gross_tonnage:e.target.value})} placeholder="15000"/></div>
              <div className={s.fld}><label>Год постройки *</label><input type="number" required min="1900" max="2100" value={form.year_built} onChange={e=>setForm({...form,year_built:e.target.value})} placeholder="2015"/></div>
            </div>
            {err && <div className={s.ferr}>{err}</div>}
            <div className={s.fbtns}>
              <button type="button" className={s.btnSec} onClick={()=>setModal(false)}>Отмена</button>
              <button type="submit" className={s.btnPri} disabled={saving}>{saving?'Сохранение...':'Сохранить'}</button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  )
}
