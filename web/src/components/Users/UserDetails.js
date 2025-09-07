import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { authService } from '../../services/auth';

const UserDetails = () => {
  const { id } = useParams();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchUser = async () => {
      try {
        setLoading(true);
        const userData = await authService.getUserById(id);
        setUser(userData);
      } catch (error) {
        setError(error.response?.data?.error || 'Failed to fetch user details');
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchUser();
    }
  }, [id]);

  if (loading) return <div style={{ textAlign: 'center', margin: '50px' }}>Loading user details...</div>;
  if (error) return <div style={{ color: 'red', textAlign: 'center', margin: '50px' }}>{error}</div>;
  if (!user) return <div style={{ textAlign: 'center', margin: '50px' }}>User not found</div>;

  return (
    <div style={{ maxWidth: '600px', margin: '50px auto', padding: '20px' }}>
      <div style={{ marginBottom: '20px' }}>
        <Link to="/users" style={{ textDecoration: 'none', color: '#007bff' }}>
          ← Back to Users List
        </Link>
      </div>

      <div style={{ 
        border: '1px solid #ddd', 
        borderRadius: '8px', 
        padding: '20px', 
        backgroundColor: '#f9f9f9',
        textAlign: 'left'
      }}>
        <h2 style={{ marginBottom: '20px', color: '#333' }}>User Details</h2>
        
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}>ID:</strong>
          <span>{user.id}</span>
        </div>
        
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}>Username:</strong>
          <span>{user.username}</span>
        </div>
        
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}>Email:</strong>
          <span>{user.email}</span>
        </div>
        
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}>Created At:</strong>
          <span>{new Date(user.created_at).toLocaleString()}</span>
        </div>
        
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}>Status:</strong>
          <span style={{ 
            color: '#28a745', 
            fontWeight: 'bold' 
          }}>
            Active
          </span>
        </div>
      </div>

      <div style={{ marginTop: '20px', textAlign: 'center' }}>
        <Link 
          to="/users" 
          style={{ 
            textDecoration: 'none', 
            padding: '10px 20px', 
            backgroundColor: '#007bff', 
            color: 'white', 
            borderRadius: '5px',
            marginRight: '10px'
          }}
        >
          Back to List
        </Link>
        
        <Link 
          to="/profile" 
          style={{ 
            textDecoration: 'none', 
            padding: '10px 20px', 
            backgroundColor: '#6c757d', 
            color: 'white', 
            borderRadius: '5px' 
          }}
        >
          My Profile
        </Link>
      </div>
    </div>
  );
};

export default UserDetails;