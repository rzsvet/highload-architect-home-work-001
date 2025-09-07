import React from 'react';
import { useAuth } from '../../context/AuthContext';

const Header = () => {
  const { user, logout } = useAuth();

  return (
    <header style={{ 
      padding: '1rem', 
      backgroundColor: '#f8f9fa', 
      borderBottom: '1px solid #dee2e6',
      marginBottom: '2rem'
    }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>User Management System</h1>
        <nav>
          {user ? (
            <div>
              <span style={{ marginRight: '1rem' }}>Welcome, {user.username}!</span>
              <a href="/profile" style={{ marginRight: '1rem' }}>Profile</a>
              <a href="/users" style={{ marginRight: '1rem' }}>Users</a>
              <button onClick={logout}>Logout</button>
            </div>
          ) : (
            <div>
              <a href="/login" style={{ marginRight: '1rem' }}>Login</a>
              <a href="/register">Register</a>
            </div>
          )}
        </nav>
      </div>
    </header>
  );
};

export default Header;