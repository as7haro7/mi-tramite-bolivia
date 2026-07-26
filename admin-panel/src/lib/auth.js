const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const getAuthTokens = () => {
  return {
    accessToken: localStorage.getItem('access_token'),
    refreshToken: localStorage.getItem('refresh_token'),
  };
};

export const setAuthTokens = (access, refresh) => {
  if (access) localStorage.setItem('access_token', access);
  if (refresh) localStorage.setItem('refresh_token', refresh);
};

export const clearAuthTokens = () => {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('usuario');
};

export const getUsuario = () => {
  const user = localStorage.getItem('usuario');
  return user ? JSON.parse(user) : null;
};

export const setUsuario = (usuario) => {
  localStorage.setItem('usuario', JSON.stringify(usuario));
};
