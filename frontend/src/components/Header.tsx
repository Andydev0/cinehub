import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Header: React.FC = () => {
  const { isLogado, logout } = useAuth();
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const toggleMenu = () => setIsMenuOpen(!isMenuOpen);

  return (
    <header className="header-gradient text-white shadow-lg sticky top-0 z-50">
      <nav className="container mx-auto px-4 py-4">
        <div className="flex justify-between items-center">
          {/* Logo */}
          <Link 
            to="/" 
            className="text-2xl font-bold text-orange hover:text-orange-light transition-colors duration-300"
          >
            🎥 CineHub
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-6">
            {isLogado ? (
              <>
                <Link 
                  to="/favoritos" 
                  className="text-gray-300 hover:text-orange transition-colors duration-300 font-medium"
                >
                  Favoritos
                </Link>
                <Link 
                  to="/recomendacoes" 
                  className="text-gray-300 hover:text-orange transition-colors duration-300 font-medium"
                >
                  Recomendações
                </Link>
                <Link 
                  to="/quiz" 
                  className="text-gray-300 hover:text-orange transition-colors duration-300 font-medium"
                >
                  Quiz
                </Link>
                <Link 
                  to="/aleatorio" 
                  className="text-gray-300 hover:text-orange transition-colors duration-300 font-medium"
                >
                  Aleatório
                </Link>
                <button 
                  onClick={logout} 
                  className="text-red-400 hover:text-red-300 transition-colors duration-300 font-medium"
                >
                  Sair
                </button>
              </>
            ) : (
              <>
                <Link 
                  to="/login" 
                  className="btn-orange-outline"
                >
                  Login
                </Link>
                <Link 
                  to="/registrar" 
                  className="btn-orange"
                >
                  Registrar
                </Link>
              </>
            )}
          </div>

          {/* Mobile Menu Button */}
          <button 
            onClick={toggleMenu}
            className="md:hidden text-orange hover:text-orange-light transition-colors duration-300"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path 
                strokeLinecap="round" 
                strokeLinejoin="round" 
                strokeWidth={2} 
                d={isMenuOpen ? "M6 18L18 6M6 6l12 12" : "M4 6h16M4 12h16M4 18h16"} 
              />
            </svg>
          </button>
        </div>

        {/* Mobile Navigation */}
        {isMenuOpen && (
          <div className="md:hidden mt-4 pb-4 border-t border-gray-700">
            <div className="flex flex-col space-y-3 pt-4">
              {isLogado ? (
                <>
                  <Link 
                    to="/favoritos" 
                    className="text-gray-300 hover:text-orange transition-colors duration-300 font-medium py-2"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    Favoritos
                  </Link>
                  <Link 
                    to="/recomendacoes" 
                    className="text-gray-300 hover:text-orange transition-colors duration-300 font-medium py-2"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    Recomendações
                  </Link>
                  <Link 
                    to="/quiz" 
                    className="text-gray-300 hover:text-orange transition-colors duration-300 font-medium py-2"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    Quiz
                  </Link>
                  <Link 
                    to="/aleatorio" 
                    className="text-gray-300 hover:text-orange transition-colors duration-300 font-medium py-2"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    Aleatório
                  </Link>
                  <button 
                    onClick={() => { logout(); setIsMenuOpen(false); }} 
                    className="text-red-400 hover:text-red-300 transition-colors duration-300 font-medium py-2 text-left"
                  >
                    Sair
                  </button>
                </>
              ) : (
                <>
                  <Link 
                    to="/login" 
                    className="btn-orange-outline text-center"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    Login
                  </Link>
                  <Link 
                    to="/registrar" 
                    className="btn-orange text-center"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    Registrar
                  </Link>
                </>
              )}
            </div>
          </div>
        )}
      </nav>
    </header>
  );
};

export default Header;