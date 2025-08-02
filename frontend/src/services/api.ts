import axios from 'axios';

// Cria a instância base do Axios.
const api = axios.create({
  baseURL: 'http://localhost:8080/v1',
});

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