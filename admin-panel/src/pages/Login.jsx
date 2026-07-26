import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../lib/api';
import { setAuthTokens, setUsuario } from '../lib/auth';
import { ShieldCheck, Loader2 } from 'lucide-react';

export function Login() {
  const [correo, setCorreo] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const data = await api.post('/admin/auth/login', { correo, password });
      setAuthTokens(data.access_token, data.refresh_token);
      setUsuario(data.usuario);
      navigate('/');
    } catch (err) {
      setError(err.data?.message || 'Error de conexión');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <div className="login-header">
          <div className="icon-wrapper">
            <ShieldCheck size={40} className="text-primary" />
          </div>
          <h1>Mi Trámite Bolivia</h1>
          <p>Panel de Administración v2</p>
        </div>

        <form onSubmit={handleLogin} className="login-form">
          {error && <div className="alert alert-error">{error}</div>}

          <div className="form-group">
            <label htmlFor="correo">Correo Electrónico</label>
            <input
              id="correo"
              type="email"
              value={correo}
              onChange={(e) => setCorreo(e.target.value)}
              placeholder="ejemplo@mitramite.bo"
              required
              disabled={loading}
              className="input-field"
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Contraseña</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
              disabled={loading}
              className="input-field"
            />
          </div>

          <button type="submit" disabled={loading} className="btn btn-primary w-full mt-4">
            {loading ? (
              <>
                <Loader2 className="spin" size={18} />
                Ingresando...
              </>
            ) : (
              'Ingresar al Sistema'
            )}
          </button>
        </form>
      </div>
    </div>
  );
}
