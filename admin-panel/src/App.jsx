import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ProtectedRoute } from './components/ProtectedRoute';
import { Layout } from './components/Layout';
import { Login } from './pages/Login';
import { Dashboard } from './pages/Dashboard';
import { TramitesList } from './pages/TramitesList';
import { InstitucionesList } from './pages/InstitucionesList';
import { TramiteDetail } from './pages/TramiteDetail';
import { OficinasList } from './pages/OficinasList';
import { ReportesList } from './pages/ReportesList';
import { UsuariosList } from './pages/UsuariosList';
import { AuditoriaList } from './pages/AuditoriaList';
import { IngestaPanel } from './pages/IngestaPanel';
import CategoriasList from './pages/CategoriasList';
import EtiquetasList from './pages/EtiquetasList';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        
        <Route element={<ProtectedRoute><Layout /></ProtectedRoute>}>
          <Route path="/" element={<Dashboard />} />
          <Route path="/tramites" element={<TramitesList />} />
          <Route path="/tramites/:id" element={<TramiteDetail />} />
          <Route path="/categorias" element={<CategoriasList />} />
          <Route path="/etiquetas" element={<EtiquetasList />} />
          <Route path="/instituciones" element={<InstitucionesList />} />
          <Route path="/oficinas" element={<OficinasList />} />
          <Route path="/reportes" element={<ReportesList />} />
          <Route path="/usuarios" element={<UsuariosList />} />
          <Route path="/auditoria" element={<AuditoriaList />} />
          <Route path="/ingesta" element={<IngestaPanel />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
