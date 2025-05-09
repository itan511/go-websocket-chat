document.addEventListener('DOMContentLoaded', function() {
      const messageForm = document.getElementById('messageForm');
      const messageInput = document.getElementById('messageInput');
      const messagesContainer = document.getElementById('messagesContainer');
      const usersList = document.getElementById('usersList');
      const leaveRoomBtn = document.getElementById('leaveRoomBtn');
      const roomName = document.getElementById('roomName');
      
      // Получаем ID комнаты из URL
      const roomId = window.location.pathname.split('/')[2];
      let socket;
      
      // Инициализация чата
      async function initChat() {
          try {
              // Получаем информацию о комнате
              const roomResponse = await fetch(`/api/rooms/${roomId}`);
              const room = await roomResponse.json();
              roomName.textContent += room.name;
              
              // Подключаемся к WebSocket
              socket = new WebSocket(`ws://${window.location.host}/ws/${roomId}`);
              
              socket.onmessage = function(event) {
                  const data = JSON.parse(event.data);
                  
                  if (data.type === 'message') {
                      addMessage(data);
                  } else if (data.type === 'users') {
                      updateUsersList(data.users);
                  }
              };
              
              // Загружаем историю сообщений
              const messagesResponse = await fetch(`/api/rooms/${roomId}/messages`);
              const messages = await messagesResponse.json();
              messages.forEach(addMessage);
              
          } catch (error) {
              console.error('Chat init error:', error);
          }
      }
      
      // Добавление сообщения
      function addMessage(message) {
          const messageElement = document.createElement('div');
          messageElement.classList.add('message', message.isCurrentUser ? 'sent' : 'received');
          messageElement.innerHTML = `
              <div class="fw-bold">${message.username}</div>
              <div>${message.content}</div>
              <div class="small text-muted">${new Date(message.timestamp).toLocaleTimeString()}</div>
          `;
          messagesContainer.appendChild(messageElement);
          messagesContainer.scrollTop = messagesContainer.scrollHeight;
      }
      
      // Обновление списка пользователей
      function updateUsersList(users) {
          usersList.innerHTML = users.map(user => `
              <li class="list-group-item">${user.username}</li>
          `).join('');
      }
      
      // Отправка сообщения
      messageForm?.addEventListener('submit', function(e) {
          e.preventDefault();
          const content = messageInput.value.trim();
          
          if (content && socket?.readyState === WebSocket.OPEN) {
              socket.send(JSON.stringify({
                  type: 'message',
                  content
              }));
              messageInput.value = '';
          }
      });
      
      // Покидание комнаты
      leaveRoomBtn?.addEventListener('click', function() {
          if (socket) socket.close();
          window.location.href = '/rooms';
      });
      
      initChat();
  });