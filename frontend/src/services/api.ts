import axios from 'axios';

// Define a URL base da API
// Usa a variável de ambiente se disponível, caso contrário usa localhost
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/v1';

// Cria a instância base do Axios.
const api = axios.create({
  baseURL: API_URL,
});

// Log para debug - remova em produção
console.log('API URL:', API_URL);

// Adiciona um interceptor que roda antes de cada requisição.
api.interceptors.request.use(
  (config) => {
    // Pega o token do localStorage.
    const token = localStorage.getItem('authToken');
    
    // Se o token existir, adiciona o header de autorização.
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    return config; // Retorna a configuração modificada.
  },
  (error) => {
    // Faz algo com o erro da requisição.
    return Promise.reject(error);
  }
);

export default api;