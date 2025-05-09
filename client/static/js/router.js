class Router {
      constructor() {
          this.routes = {
              '/': 'auth.html',
              '/auth': 'auth.html',
              '/rooms': 'rooms.html',
              '/chat/:id': 'chat.html'
          };
          this.init();
      }
  
      init() {
          window.addEventListener('popstate', () => this.route());
          document.addEventListener('DOMContentLoaded', () => this.route());
          document.addEventListener('click', (e) => {
              if (e.target.closest('[data-link]')) {
                  e.preventDefault();
                  this.navigate(e.target.closest('[data-link]').href);
              }
          });
      }
  
      async route() {
          const path = window.location.pathname;
          const matchingRoute = Object.keys(this.routes).find(route => {
              const routeParts = route.split('/');
              const pathParts = path.split('/');
              
              if (routeParts.length !== pathParts.length) return false;
              
              return routeParts.every((part, i) => 
                  part.startsWith(':') || part === pathParts[i]
              );
          });
  
          const template = matchingRoute ? this.routes[matchingRoute] : 'auth.html';
          await this.loadTemplate(template);
          this.loadController(path);
      }
  
      async loadTemplate(template) {
          try {
              const html = await fetch(template).then(res => res.text());
              document.getElementById('app').innerHTML = html;
          } catch (error) {
              console.error('Failed to load template:', error);
          }
      }
  
      loadController(path) {
          const scriptMap = {
              '/auth': '/static/js/auth.js',
              '/rooms': '/static/js/rooms.js',
              '/chat': '/static/js/chat.js'
          };
  
          const scriptPath = Object.keys(scriptMap).find(key => path.startsWith(key));
          if (scriptPath) {
              this.loadScript(scriptMap[scriptPath]);
          }
      }
  
      loadScript(src) {
          if (document.querySelector(`script[src="${src}"]`)) return;
          
          const script = document.createElement('script');
          script.src = src;
          document.body.appendChild(script);
      }
  
      navigate(path) {
          window.history.pushState({}, '', path);
          this.route();
      }
  }
  
  new Router();