const BASE = '/api'
const TOKEN_KEY = 'harbor_token'
const USER_KEY  = 'harbor_user'


export const saveAuth  = (token, user) => { localStorage.setItem(TOKEN_KEY, token); localStorage.setItem(USER_KEY, JSON.stringify(user)) }
export const clearAuth = () => { localStorage.removeItem(TOKEN_KEY); localStorage.removeItem(USER_KEY) }
export const getToken  = () => localStorage.getItem(TOKEN_KEY)
export const getUser   = () => { const r = localStorage.getItem(USER_KEY); return r ? JSON.parse(r) : null }


const req = async (method, path, body) => {
  const headers = { 'Content-Type': 'application/json' }
  const token = getToken()
  if (token) headers['Authorization'] = `Bearer ${token}`
  const res = await fetch(BASE + path, { method, headers, body: body ? JSON.stringify(body) : undefined })
  if (res.status === 401) { clearAuth(); window.location.href = '/login'; throw new Error('Сессия истекла') }
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || `Ошибка ${res.status}`)
  return data
}


export const login   = (email, password) => req('POST', '/auth/login', { email, password })
export const getMe   = () => req('GET', '/me')


export const getUsers    = ()       => req('GET',    '/users')
export const getUserById = (id)     => req('GET',    `/users/${id}`)
export const createUser  = (data)   => req('POST',   '/users',    data)
export const updateUser  = (id, d)  => req('PUT',    `/users/${id}`, d)
export const deleteUser  = (id)     => req('DELETE', `/users/${id}`)


export const getShips    = ()       => req('GET',    '/ships')
export const getShipById = (id)     => req('GET',    `/ships/${id}`)
export const createShip  = (data)   => req('POST',   '/ships',    data)
export const updateShip  = (id, d)  => req('PUT',    `/ships/${id}`, d)
export const deleteShip  = (id)     => req('DELETE', `/ships/${id}`)


export const getVisits    = ()       => req('GET',    '/visits')
export const getVisitById = (id)     => req('GET',    `/visits/${id}`)
export const createVisit  = (data)   => req('POST',   '/visits',    data)
export const updateVisit  = (id, d)  => req('PUT',    `/visits/${id}`, d)
export const deleteVisit  = (id)     => req('DELETE', `/visits/${id}`)
