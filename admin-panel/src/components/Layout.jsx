import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { clearAuthTokens, getUsuario } from '../lib/auth';
import { api } from '../lib/api';
import { 
  LayoutDashboard, 
  Files, 
  Building2, 
  LogOut, 
  Menu,
  X,
  UserCircle,
  MapPin,
  MessageSquareWarning,
  Users,
  History,
  Database,
  FolderTree,
  Tags
} from 'lucide-react';
import { useState, useEffect } from 'react';

export function Layout() {
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [usuario, setUsuario] = useState(null);

  useEffect(() => {
    // Intentar obtener de localStorage primero
    const stored = getUsuario();
    if (stored) {
      setUsuario(stored);
    }
  }, []);

  const handleLogout = async () => {
    try {
      await api.post('/admin/auth/logout', {});
    } catch (e) {
      // Ignore errors on logout
    } finally {
      clearAuthTokens();
      navigate('/login');
    }
  };

  const navItems = [
    { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
    { to: '/tramites', icon: Files, label: 'Trámites' },
    { to: '/categorias', icon: FolderTree, label: 'Categorías' },
    { to: '/etiquetas', icon: Tags, label: 'Etiquetas' },
    { to: '/instituciones', icon: Building2, label: 'Instituciones' },
    { to: '/oficinas', icon: MapPin, label: 'Puntos de Atención' },
    { to: '/ingesta', icon: Database, label: 'Ingesta IA' },
    { to: '/reportes', icon: MessageSquareWarning, label: 'Reportes Ciudadanos' },
    { to: '/usuarios', icon: Users, label: 'Usuarios' },
    { to: '/auditoria', icon: History, label: 'Auditoría' },
  ];

  const SidebarContent = () => (
    <>
      <div className="sidebar-header">
        <h2>Mi Trámite</h2>
        <span className="badge-env">Admin v2</span>
      </div>
      
      <nav className="sidebar-nav">
        {navItems.map((item) => (
          <NavLink 
            key={item.to} 
            to={item.to} 
            className={({isActive}) => isActive ? 'nav-item active' : 'nav-item'}
            onClick={() => setSidebarOpen(false)}
          >
            <item.icon size={20} />
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>

      <div className="sidebar-footer">
        <div className="user-profile">
          <UserCircle size={32} />
          <div className="user-info">
            <span className="user-name">{usuario?.nombre || 'Admin'}</span>
            <span className="user-role">{usuario?.roles?.[0] || 'Editor'}</span>
          </div>
        </div>
        <button onClick={handleLogout} className="logout-btn">
          <LogOut size={18} />
          <span>Cerrar Sesión</span>
        </button>
      </div>
    </>
  );

  return (
    <div className="app-layout">
      {/* Mobile Header */}
      <div className="mobile-header">
        <h2>Mi Trámite</h2>
        <button className="icon-btn" onClick={() => setSidebarOpen(true)}>
          <Menu size={24} />
        </button>
      </div>

      {/* Sidebar Overlay */}
      {sidebarOpen && (
        <div className="sidebar-overlay" onClick={() => setSidebarOpen(false)}></div>
      )}

      {/* Sidebar */}
      <aside className={`sidebar ${sidebarOpen ? 'open' : ''}`}>
        {sidebarOpen && (
          <button className="icon-btn close-sidebar" onClick={() => setSidebarOpen(false)}>
            <X size={24} />
          </button>
        )}
        <SidebarContent />
      </aside>

      {/* Main Content */}
      <main className="main-content">
        <div className="page-container">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
