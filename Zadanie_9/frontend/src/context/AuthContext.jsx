import { createContext, useState, useContext, useEffect, useMemo, useCallback } from 'react';
import api from '../api';

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [loading, setLoading] = useState(!!localStorage.getItem('token'));

  useEffect(() => {
    if (token) {
      api.get('/auth/me')
        .then(res => {
          setUser(res.data);
          setLoading(false);
        })
        .catch(() => {
          localStorage.removeItem('token');
          setToken(null);
          setUser(null);
          setLoading(false);
        });
    }
  }, [token]);

  const login = useCallback(async (email, password) => {
    const res = await api.post('/auth/login', { email, password });
    const { token: newToken, user: userData } = res.data;
    localStorage.setItem('token', newToken);
    setToken(newToken);
    setUser(userData);
    return userData;
  }, []);

  const register = useCallback(async (email, password, name) => {
    const res = await api.post('/auth/register', { email, password, name });
    const { token: newToken, user: userData } = res.data;
    localStorage.setItem('token', newToken);
    setToken(newToken);
    setUser(userData);
    return userData;
  }, []);

  const saveToken = useCallback((newToken) => {
    localStorage.setItem('token', newToken);
    setToken(newToken);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
  }, []);

  const value = useMemo(() => ({
    user,
    token,
    loading,
    login,
    register,
    saveToken,
    logout,
  }), [user, token, loading, login, register, saveToken, logout]);

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
