import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { authService } from '../../services/auth';

const UserProfile = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const profileData = await authService.getProfile();
        setProfile(profileData);
      } catch (error) {
        setError(error.response?.data?.error || 'Failed to fetch profile');
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div style={{ color: 'red' }}>{error}</div>;

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
        <h2 style={{ marginBottom: '20px', color: '#333' }}>User Profile</h2>

        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> ID:</strong>
          <span>{profile.id}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> Username:</strong>
          <span>{profile.username}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> Email:</strong>
          <span>{profile.email}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> FirstName:</strong>
          <span>{profile.first_name}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> LastName:</strong>
          <span>{profile.last_name}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> BirthDate:</strong>
          <span>{new Date(profile.birth_date).toLocaleString()}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> Gender:</strong>
          <span>{profile.gender}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> Interests:</strong>
          <span>{profile.interests}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> City:</strong>
          <span>{profile.city}</span>
        </div>
        <div style={{ marginBottom: '15px' }}>
          <strong style={{ display: 'inline-block', width: '100px' }}> Created At:</strong>
          <span>{new Date(profile.created_at).toLocaleString()}</span>
        </div>
      </div>
    </div>
  );
};

export default UserProfile;