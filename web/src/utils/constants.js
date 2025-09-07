// Синхронная загрузка конфига
let config = {};

try {
  // Пытаемся загрузить конфиг синхронно (не рекомендуется для production)
  const xhr = new XMLHttpRequest();
  xhr.open('GET', '/config.json', false); // sync request
  xhr.send();
  
  if (xhr.status === 200) {
    config = JSON.parse(xhr.responseText);
  }
} catch (error) {
  console.warn('Failed to load config.json, using default values:', error);
}

export const API_BASE_URL = config.API_BASE_URL || 'http://localhost:8080/api/v1';
export const FRONTEND_BASE_URL = config.FRONTEND_BASE_URL || 'http://localhost:3000';