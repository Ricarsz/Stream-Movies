import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth
export const registerUser = (data) => api.post('/users/register', data);
export const loginUser = (data) => api.post('/users/login', data);

// Movies
export const getMovies = () => api.get('/movies');
export const getMovie = (imdbId) => api.get(`/movies/${imdbId}`);
export const addMovie = (data) => api.post('/movies', data);
export const updateMovie = (imdbId, data) => api.put(`/movies/${imdbId}`, data);
export const deleteMovie = (imdbId) => api.delete(`/movies/${imdbId}`);

export default api;
