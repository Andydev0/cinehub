import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegistroPage from './pages/RegistroPage';
import ResultadosBuscaPage from './pages/ResultadosBuscaPage';
import FavoritosPage from './pages/FavoritosPage';
import RotaProtegida from './components/RotaProtegida';
import RecomendacoesPage from './pages/RecomendacoesPage';
import QuizPage from './pages/QuizPage';
import GeradorAleatorioPage from './pages/GeradorAleatorioPage';
import PaginaDetalhesFilme from './pages/PaginaDetalhesFilme';

function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen gradient-bg flex flex-col">
        <Header />
        <main className="text-white flex-1">
          <Routes>
            {/* Rotas públicas */}
            <Route path="/" element={<HomePage />} />
            <Route path="/buscar" element={<ResultadosBuscaPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/registrar" element={<RegistroPage />} />
            <Route path="/aleatorio" element={<GeradorAleatorioPage />} />
            <Route path="/filme/:id" element={<PaginaDetalhesFilme />} />

            {/* Rotas protegidas */}
            <Route element={<RotaProtegida />}>
              <Route path="/favoritos" element={<FavoritosPage />} />
              <Route path="/recomendacoes" element={<RecomendacoesPage />} />
              <Route path="/quiz" element={<QuizPage />} />
            </Route>
          </Routes>
        </main>
        <Footer />
      </div>
    </BrowserRouter>
  );
}

export default App;