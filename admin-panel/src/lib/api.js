import { getAuthTokens, setAuthTokens, clearAuthTokens } from './auth';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

class ApiError extends Error {
  constructor(status, data) {
    super(data.message || 'Error en la petición');
    this.status = status;
    this.data = data;
  }
}

async function request(endpoint, options = {}) {
  const { accessToken } = getAuthTokens();
  
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };

  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`;
  }

  const config = {
    ...options,
    headers,
  };

  const url = endpoint.startsWith('http') ? endpoint : `${API_BASE}${endpoint}`;

  let response = await fetch(url, config);

  // Intentar refrescar token si es 401
  if (response.status === 401 && accessToken) {
    const { refreshToken } = getAuthTokens();
    if (refreshToken) {
      try {
        const refreshResponse = await fetch(`${API_BASE}/admin/auth/refresh`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refresh_token: refreshToken })
        });
        
        if (refreshResponse.ok) {
          const refreshData = await refreshResponse.json();
          setAuthTokens(refreshData.access_token, null);
          
          // Reintentar la petición original con el nuevo token
          config.headers['Authorization'] = `Bearer ${refreshData.access_token}`;
          response = await fetch(url, config);
        } else {
          // Refresh fallido, forzar logout
          clearAuthTokens();
          window.location.href = '/login';
          throw new Error('Sesión expirada');
        }
      } catch (err) {
        clearAuthTokens();
        window.location.href = '/login';
        throw err;
      }
    } else {
      clearAuthTokens();
      window.location.href = '/login';
    }
  }

  const isJson = response.headers.get('content-type')?.includes('application/json');
  const data = isJson ? await response.json() : await response.text();

  if (!response.ok) {
    throw new ApiError(response.status, data);
  }

  return data;
}

export const api = {
  get: (endpoint, options) => request(endpoint, { method: 'GET', ...options }),
  post: (endpoint, body, options) => request(endpoint, { method: 'POST', body: JSON.stringify(body), ...options }),
  put: (endpoint, body, options) => request(endpoint, { method: 'PUT', body: JSON.stringify(body), ...options }),
  delete: (endpoint, options) => request(endpoint, { method: 'DELETE', ...options }),
};
