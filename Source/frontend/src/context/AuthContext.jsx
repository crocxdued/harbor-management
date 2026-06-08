import { createContext, useContext, useState, useEffect } from 'react'
import { getToken, getUser, saveAuth, clearAuth, login as apiLogin } from '../api/api'

const Ctx = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (getToken() && getUser()) setUser(getUser())
    setLoading(false)
  }, [])

  const login = async (email, password) => {
    const r = await apiLogin(email, password)
    saveAuth(r.token, r.user)
    setUser(r.user)
    return r
  }

  const logout = () => { clearAuth(); setUser(null) }
  const hasRole = (...roles) => user && roles.includes(user.role)

  return <Ctx.Provider value={{ user, loading, login, logout, hasRole }}>{children}</Ctx.Provider>
}

export const useAuth = () => useContext(Ctx)
