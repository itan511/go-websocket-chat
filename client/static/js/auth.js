document.addEventListener('DOMContentLoaded', function() {
      // Логика авторизации
      const loginForm = document.getElementById('loginForm');
      const registerForm = document.getElementById('registerForm');
      
      loginForm?.addEventListener('submit', async function(e) {
          e.preventDefault();
          const email = this.querySelector('input[type="email"]').value;
          const password = this.querySelector('input[type="password"]').value;
          
          try {
              const response = await fetch('/api/auth/login', {
                  method: 'POST',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({ email, password })
              });
              
              if (response.ok) {
                  window.location.href = '/rooms';
              } else {
                  alert('Ошибка авторизации');
              }
          } catch (error) {
              console.error('Login error:', error);
          }
      });
      
      registerForm?.addEventListener('submit', async function(e) {
          e.preventDefault();
          const username = this.querySelector('input[type="text"]').value;
          const email = this.querySelector('input[type="email"]').value;
          const password = this.querySelector('input[type="password"]').value;
          
          try {
              const response = await fetch('/api/auth/register', {
                  method: 'POST',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({ username, email, password })
              });
              
              if (response.ok) {
                  alert('Регистрация успешна! Теперь войдите');
                  document.querySelector('.nav-link[href="#login"]').click();
              } else {
                  alert('Ошибка регистрации');
              }
          } catch (error) {
              console.error('Registration error:', error);
          }
      });
  });