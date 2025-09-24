import api from './api';

export const authService = {
  register: async (userData) => {
    const response = await api.post('/user/register', userData);
    return response.data;
  },

  login: async (credentials) => {
    const response = await api.post('/login', credentials);
    return response.data;
  },

  getProfile: async () => {
    const response = await api.get('/profile');
    return response.data;
  },

  getAllUsers: async () => {
    const response = await api.get('/users');
    return response.data;
  },

  getUserById: async (id) => {
    const response = await api.get(`/user/get/${id}`);
    return response.data;
  },

  searchUsers: async (firstName, lastName, page = 1, pageSize = 20) => {
    const response = await api.get('/user/search', {
      params: {
        first_name: firstName,
        last_name: lastName,
        page: page,
        page_size: pageSize
      }
    });
    return response.data;
  }
};