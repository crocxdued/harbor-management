import s from './Modal.module.css'

export default function Modal({ title, onClose, children }) {
  return (
    <div className={s.overlay} onClick={onClose}>
      <div className={s.box} onClick={e => e.stopPropagation()}>
        <div className={s.head}>
          <h2 className={s.title}>{title}</h2>
          <button className={s.x} onClick={onClose}>✕</button>
        </div>
        <div className={s.body}>{children}</div>
      </div>
    </div>
  )
}
