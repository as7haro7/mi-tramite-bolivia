import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { Users, Plus, Shield, CheckCircle2, XCircle } from 'lucide-react';
import { Modal } from '../components/Modal';

export function UsuariosList() {
  const [usuarios, setUsuarios] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form, setForm] = useState({ 
    nombre: '', 
    correo: '', 
    password: '', 
    rol_codigo: 'editor', // Default: Editor 
    estado: 'activo' 
  });
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    cargarDatos();
  }, []);

  const cargarDatos = async () => {
    try {
      setLoading(true);
      const resp = await api.get('/admin/usuarios');
      setUsuarios(resp.datos || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setFormLoading(true);
    try {
      if (form.id) {
        await api.put(`/admin/usuarios/${form.id}`, {
          nombre: form.nombre,
          rol_codigo: form.rol_codigo
        });
      } else {
        await api.post('/admin/usuarios', {
          ...form
        });
      }
      setIsModalOpen(false);
      setForm({ id: null, nombre: '', correo: '', password: '', rol_codigo: 'editor', estado: 'activo' });
      cargarDatos();
    } catch (err) {
      alert(err.data?.message || 'Error al guardar usuario');
    } finally {
      setFormLoading(false);
    }
  };

  const openCreate = () => {
    setForm({ id: null, nombre: '', correo: '', password: '', rol_codigo: 'editor', estado: 'activo' });
    setIsModalOpen(true);
  };

  const openEdit = (u) => {
    setForm({
      id: u.id,
      nombre: u.nombre,
      correo: u.correo, // Not editable but used for display
      password: '', // Hidden in edit
      rol_codigo: u.roles?.[0] || 'editor',
      estado: u.estado
    });
    setIsModalOpen(true);
  };

  const handleToggleEstado = async (u) => {
    const nuevoEstado = u.estado === 'activo' ? 'bloqueado' : 'activo';
    try {
      await api.put(`/admin/usuarios/${u.id}`, { estado: nuevoEstado });
      cargarDatos();
    } catch (err) {
      alert('Error al cambiar estado');
    }
  };

  if (loading && usuarios.length === 0) return <div className="page-loader"><Users className="spin" size={32} /></div>;
  if (error) return <div className="alert alert-error">{error}</div>;

  return (
    <div className="usuarios-list">
      <header className="page-header flex justify-between items-center">
        <div>
          <h1 className="page-title">Usuarios Administrativos</h1>
          <p className="page-subtitle">Gestión de editores, administradores y permisos</p>
        </div>
        <button className="btn btn-primary" onClick={openCreate}>
          <Plus size={18} /> Nuevo Usuario
        </button>
      </header>

      <div className="card mt-6">
        <div className="table-responsive">
          <table className="table">
            <thead>
              <tr>
                <th>Nombre</th>
                <th>Correo Electrónico</th>
                <th>Roles (Globales)</th>
                <th>Estado</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {usuarios.length === 0 ? (
                <tr>
                  <td colSpan="5" className="text-center py-8 text-neutral-500">
                    No hay usuarios registrados.
                  </td>
                </tr>
              ) : (
                usuarios.map(u => (
                  <tr key={u.id}>
                    <td className="font-bold">{u.nombre}</td>
                    <td className="text-sm">{u.correo}</td>
                    <td>
                      {u.roles?.map(r => (
                        <span key={r} className="badge badge-primary mr-2 mb-1 inline-flex items-center gap-1">
                          <Shield size={12} /> {r}
                        </span>
                      ))}
                    </td>
                    <td>
                      <span className={`badge ${u.estado === 'activo' ? 'badge-success' : 'badge-error'}`}>
                        {u.estado.toUpperCase()}
                      </span>
                    </td>
                    <td>
                      <button 
                        className="btn btn-sm btn-outline mr-2" 
                        onClick={() => openEdit(u)}
                      >
                        Editar
                      </button>
                      <button 
                        className="btn btn-sm btn-outline" 
                        onClick={() => handleToggleEstado(u)}
                      >
                        {u.estado === 'activo' ? 'Bloquear' : 'Activar'}
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={form.id ? 'Editar Usuario' : 'Nuevo Usuario'}>
        <form onSubmit={handleSubmit} className="flex-col">
          <div className="form-group">
            <label>Nombre Completo</label>
            <input required className="input-field" value={form.nombre} onChange={e => setForm({...form, nombre: e.target.value})} placeholder="Ej. Juan Pérez" />
          </div>
          
          <div className="form-group">
            <label>Correo Electrónico</label>
            <input required type="email" className="input-field" value={form.correo} onChange={e => setForm({...form, correo: e.target.value})} placeholder="juan@mitramite.bo" disabled={!!form.id} />
          </div>

          {!form.id && (
            <div className="form-group">
              <label>Contraseña Temporal</label>
              <input required type="password" className="input-field" value={form.password} onChange={e => setForm({...form, password: e.target.value})} placeholder="Mínimo 8 caracteres" minLength={8} />
            </div>
          )}

          <div className="form-group">
            <label>Rol Global</label>
            <select required className="input-field" value={form.rol_codigo} onChange={e => setForm({...form, rol_codigo: e.target.value})}>
              <option value="superadmin">Superadmin (Acceso total)</option>
              <option value="editor">Editor (Crear y editar trámites)</option>
              <option value="revisor">Revisor (Aprobar trámites)</option>
              <option value="analista">Analista (Solo lectura)</option>
            </select>
          </div>

          <div className="flex justify-between mt-6">
            <button type="button" className="btn btn-outline" onClick={() => setIsModalOpen(false)}>Cancelar</button>
            <button type="submit" className="btn btn-primary" disabled={formLoading}>
              {formLoading ? 'Guardando...' : (form.id ? 'Guardar Cambios' : 'Crear Usuario')}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
