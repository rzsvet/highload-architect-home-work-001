import React, { useState } from 'react';
import { authService } from '../../services/auth';
import './UserSearch.css';

const UserSearch = () => {
  const [searchData, setSearchData] = useState({
    first_name: '',
    last_name: ''
  });
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [pagination, setPagination] = useState({
    page: 1,
    page_size: 20,
    total: 0,
    total_pages: 0
  });

  const handleChange = (e) => {
    setSearchData({
      ...searchData,
      [e.target.name]: e.target.value
    });
  };

  const handleSearch = async (page = 1) => {
    if (searchData.first_name.length < 2 || searchData.last_name.length < 2) {
      setError('Имя и фамилия должны содержать минимум 2 символа');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await authService.searchUsers(
        searchData.first_name,
        searchData.last_name,
        page,
        pagination.page_size
      );
      
      setResults(response.users || response.items || response.data || []);
      setPagination(prev => ({
        ...prev,
        page: page,
        total: response.total || response.total_count || 0,
        total_pages: response.total_pages || Math.ceil((response.total || 0) / pagination.page_size)
      }));
    } catch (error) {
      setError(error.response?.data?.error || 'Ошибка при поиске пользователей');
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    handleSearch(1);
  };

  const handlePageChange = (newPage) => {
    if (newPage >= 1 && newPage <= pagination.total_pages) {
      handleSearch(newPage);
    }
  };

  const handlePageSizeChange = (e) => {
    const newSize = parseInt(e.target.value);
    setPagination(prev => ({ ...prev, page_size: newSize }));
    // Перезагружаем поиск с новой страницей
    setTimeout(() => handleSearch(1), 100);
  };

  return (
    <div className="user-search-container">
      <div className="user-search-card">
        <h2 className="search-title">Поиск пользователей</h2>
        <p className="search-subtitle">Найдите пользователей по имени и фамилии</p>
        
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit} className="search-form">
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="first_name" className="form-label">Имя *</label>
              <input
                id="first_name"
                type="text"
                name="first_name"
                value={searchData.first_name}
                onChange={handleChange}
                required
                minLength={2}
                className="form-input"
                placeholder="Введите имя (минимум 2 символа)"
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="last_name" className="form-label">Фамилия *</label>
              <input
                id="last_name"
                type="text"
                name="last_name"
                value={searchData.last_name}
                onChange={handleChange}
                required
                minLength={2}
                className="form-input"
                placeholder="Введите фамилию (минимум 2 символа)"
              />
            </div>
          </div>

          <div className="search-controls">
            <div className="page-size-selector">
              <label htmlFor="page_size" className="form-label">Пользователей на странице:</label>
              <select
                id="page_size"
                value={pagination.page_size}
                onChange={handlePageSizeChange}
                className="form-select"
              >
                <option value={10}>10</option>
                <option value={20}>20</option>
                <option value={50}>50</option>
                <option value={100}>100</option>
              </select>
            </div>
            
            <button 
              type="submit" 
              disabled={loading || searchData.first_name.length < 2 || searchData.last_name.length < 2}
              className="search-button"
            >
              {loading ? (
                <span className="button-loading">
                  <span className="spinner"></span>
                  Поиск...
                </span>
              ) : (
                'Найти пользователей'
              )}
            </button>
          </div>
        </form>

        {/* Результаты поиска */}
        {results.length > 0 && (
          <div className="search-results">
            <h3 className="results-title">
              Найдено пользователей: {pagination.total}
              {pagination.total_pages > 1 && ` (Страница ${pagination.page} из ${pagination.total_pages})`}
            </h3>
            
            <div className="users-table-container">
              <table className="users-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Имя</th>
                    <th>Фамилия</th>
                    <th>Username</th>
                    <th>Email</th>
                    <th>Город</th>
                    <th>Дата регистрации</th>
                    <th>Действия</th>
                  </tr>
                </thead>
                <tbody>
                  {results.map((user) => (
                    <tr key={user.id}>
                      <td>{user.id}</td>
                      <td>{user.first_name}</td>
                      <td>{user.last_name}</td>
                      <td>{user.username}</td>
                      <td>{user.email}</td>
                      <td>{user.city || '-'}</td>
                      <td>{new Date(user.created_at).toLocaleDateString('ru-RU')}</td>
                      <td>
                        <a 
                          href={`/user/${user.id}`} 
                          className="view-link"
                          title="Посмотреть профиль"
                        >
                          👁️
                        </a>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Пагинация */}
            {pagination.total_pages > 1 && (
              <div className="pagination">
                <button
                  onClick={() => handlePageChange(1)}
                  disabled={pagination.page === 1}
                  className="pagination-button"
                >
                  ⏮️
                </button>
                
                <button
                  onClick={() => handlePageChange(pagination.page - 1)}
                  disabled={pagination.page === 1}
                  className="pagination-button"
                >
                  ◀️
                </button>

                {Array.from({ length: Math.min(5, pagination.total_pages) }, (_, i) => {
                  let pageNum;
                  if (pagination.total_pages <= 5) {
                    pageNum = i + 1;
                  } else if (pagination.page <= 3) {
                    pageNum = i + 1;
                  } else if (pagination.page >= pagination.total_pages - 2) {
                    pageNum = pagination.total_pages - 4 + i;
                  } else {
                    pageNum = pagination.page - 2 + i;
                  }

                  return (
                    <button
                      key={pageNum}
                      onClick={() => handlePageChange(pageNum)}
                      className={`pagination-button ${pagination.page === pageNum ? 'active' : ''}`}
                    >
                      {pageNum}
                    </button>
                  );
                })}

                <button
                  onClick={() => handlePageChange(pagination.page + 1)}
                  disabled={pagination.page === pagination.total_pages}
                  className="pagination-button"
                >
                  ▶️
                </button>
                
                <button
                  onClick={() => handlePageChange(pagination.total_pages)}
                  disabled={pagination.page === pagination.total_pages}
                  className="pagination-button"
                >
                  ⏭️
                </button>
              </div>
            )}
          </div>
        )}

        {results.length === 0 && !loading && searchData.first_name && searchData.last_name && (
          <div className="no-results">
            <p>Пользователи не найдены. Попробуйте изменить параметры поиска.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default UserSearch;