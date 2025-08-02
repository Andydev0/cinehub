import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';

const RotaProtegida: React.FC = () => {
  const token = localStorage.getItem('authToken');

  // Se não há token, redireciona para a página de login.
  if (!token) {
    return <Navigate to="/login" replace />;
  }

  // Se há token, renderiza o componente da rota (ex: FavoritosPage).
  return <Outlet />;
};

export default RotaProtegida;