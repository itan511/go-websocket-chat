document.addEventListener('DOMContentLoaded', function() {
      const roomsList = document.getElementById('roomsList');
      const createRoomForm = document.getElementById('createRoomForm');
      const logoutBtn = document.getElementById('logoutBtn');
      
      // Загрузка комнат
      async function loadRooms() {
          try {
              const response = await fetch('/api/rooms');
              const rooms = await response.json();
              
              roomsList.innerHTML = rooms.map(room => `
                  <a href="/chat/${room.id}" class="list-group-item list-group-item-action" data-link>
                      ${room.name}
                      <span class="badge bg-primary rounded-pill float-end">${room.users_count}</span>
                  </a>
              `).join('');
          } catch (error) {
              console.error('Failed to load rooms:', error);
          }
      }
      
      // Создание комнаты
      createRoomForm?.addEventListener('submit', async function(e) {
          e.preventDefault();
          const name = this.querySelector('input').value;
          
          try {
              const response = await fetch('/api/rooms', {
                  method: 'POST',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({ name })
              });
              
              if (response.ok) {
                  this.reset();
                  loadRooms();
              }
          } catch (error) {
              console.error('Failed to create room:', error);
          }
      });
      
      // Выход
      logoutBtn?.addEventListener('click', async function() {
          await fetch('/api/auth/logout', { method: 'POST' });
          window.location.href = '/auth';
      });
      
      loadRooms();
  });