import { Navigate, useLocation } from 'react-router-dom';
import { getAuthTokens } from '../lib/auth';

export function ProtectedRoute({ children }) {
  const { accessToken } = getAuthTokens();
  const location = useLocation();

  if (!accessToken) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return children;
}
