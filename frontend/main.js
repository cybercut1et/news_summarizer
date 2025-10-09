// --- frontend/main.js (cleaned) ---
// UI helpers and server interaction
const loginBtn = document.getElementById('loginBtn');
const loginModal = document.getElementById('loginModal');
const loginForm = document.getElementById('loginForm');
const overlay = document.getElementById('overlay');
const usernameInput = document.getElementById('username');
const passwordInput = document.getElementById('password');
// frontend/main.js — compact and valid
document.addEventListener('DOMContentLoaded', () => {
  const loginBtn = document.getElementById('loginBtn');
  const loginModal = document.getElementById('loginModal');
  const loginForm = document.getElementById('loginForm');
  const overlay = document.getElementById('overlay');
  const usernameInput = document.getElementById('username');
  const cancel = document.getElementById('cancel');

  const plusBtn = document.getElementById('plus');
  const minusBtn = document.getElementById('minus');
  const containerUrl = document.getElementById('container-url');
  const apply = document.getElementById('apply');

  function openModal() { if (!loginModal) return; loginModal.style.display = 'flex'; overlay.style.display = 'block'; document.body.style.overflow = 'hidden'; }
  function closeModal() { if (!loginModal) return; loginModal.style.display = 'none'; overlay.style.display = 'none'; document.body.style.overflow = ''; if (loginForm) loginForm.reset(); }

  if (loginBtn) loginBtn.addEventListener('click', openModal);
  if (cancel) cancel.addEventListener('click', closeModal);
  if (overlay) overlay.addEventListener('click', closeModal);

  if (loginForm) loginForm.addEventListener('submit', (e) => { e.preventDefault(); const username = usernameInput ? usernameInput.value : ''; localStorage.setItem('username', username); updateAuthButton(); closeModal(); });

  function updateAuthButton() { const b = document.getElementById('loginBtn'); const u = localStorage.getItem('username'); if (!b) return; if (u) { b.innerHTML = '<img src="../frontend/icons/free-icon-font-user-3917711.png" style="width:24px;">'; b.className = 'avatar-btn'; } else { b.innerHTML = 'Вход'; b.className = 'login-btn'; } }

  if (plusBtn) plusBtn.addEventListener('click', () => { const d = document.createElement('div'); d.className = 'url-box'; d.innerHTML = '<input type="url" placeholder="ссылка на ресурс">'; if (containerUrl) containerUrl.appendChild(d); });
  if (minusBtn) minusBtn.addEventListener('click', () => { const windows = document.querySelectorAll('.url-box'); if (windows.length > 1) { const last = windows[windows.length-1]; last.remove(); } });

  document.querySelectorAll('.categories-list button').forEach(b => b.addEventListener('click', () => b.classList.toggle('active')));

  const range = document.querySelector('.range-input'); if (range) range.addEventListener('input', (e) => { e.target.parentNode.parentNode.style.setProperty('--value', e.target.value); if (e.target.nextElementSibling) e.target.nextElementSibling.value = e.target.value; });

  function showLoading(on) { const w = document.querySelector('.welcome'); if (!w) return; if (on) { w.innerHTML = '<div class="loading"><img src="../frontend/icons/free-icon-old-scroll-17472284.png" style="width:110px;"><div class="loading-text">Идет поиск новостей...</div></div>'; w.style.display = 'block'; } else { w.style.display = 'none'; } }

  function createNewsBlock(item) {
    const block = document.createElement('div');
    block.className = 'class-block-news';
    // Показываем полный текст новости (content). Если content пустой, используем summary.
    const bodyText = item.content && item.content.trim() !== '' ? item.content : (item.summary || item.text || '');
    block.innerHTML = `${item.category ? `<div class="cat">${item.category}</div>` : ''}<div class="text">${bodyText}</div>${item.url ? `<a class="link" href="${item.url}" target="_blank">Перейти в источник</a>` : ''}`;
    const container = document.querySelector('.news');
    if (container) container.appendChild(block);
  }

  function createMultipleNewsBlocks(arr) { const container = document.querySelector('.news'); if (!container) return; container.innerHTML = ''; arr.forEach(createNewsBlock); }

  function collectFiltersData() { const urls = []; if (containerUrl) containerUrl.querySelectorAll('input[type="url"]').forEach(i => { if (i.value.trim()) urls.push(i.value.trim()); }); const cats = []; document.querySelectorAll('.categories-list button.active').forEach(b => cats.push(b.textContent.trim())); const periodEl = document.getElementById('tailmetr'); const period = periodEl ? parseInt(periodEl.value) : 24; return { urls, categories: cats, period }; }

  function loadNews() {
    const filters = collectFiltersData(); 
    const welcome = document.querySelector('.welcome'); 
    if (welcome) welcome.style.display = 'none'; 
    showLoading(true);
    
    const params = new URLSearchParams();
    if (filters.categories && filters.categories.length) params.set('categories', filters.categories.join(','));
    const sentencesSelect = document.getElementById('sentences-select');
    const sentences = (sentencesSelect && sentencesSelect.value) ? sentencesSelect.value : '2';
    params.set('sentences', sentences);

    fetch('/api/news?' + params.toString())
      .then(r => { if (!r.ok) throw new Error('news fetch failed'); return r.json(); })
      .then(data => { 
        // Проверяем, загружаются ли данные
        if (data.status === 'loading') {
          // Показываем статус пайплайна и проверяем снова
          updateLoadingStatus(data.pipeline_status || 'Обработка данных...');
          setTimeout(() => loadNews(), 2000); // Проверяем каждые 2 секунды
          return;
        }
        
        showLoading(false); 
        if (Array.isArray(data) && data.length > 0) {
          createMultipleNewsBlocks(data);
        } else {
          const w = document.querySelector('.welcome');
          if (w) {
            w.innerHTML = 'Новостей не найдено. Данные обрабатываются...';
            w.style.display = 'block';
          }
        }
      })
      .catch(err => { showLoading(false); console.error(err); alert('Ошибка при загрузке новостей'); });
  }

  function updateLoadingStatus(status) {
    const loadingText = document.querySelector('.loading-text');
    if (loadingText) {
      loadingText.textContent = status;
    }
  }

  function checkPipelineStatus() {
    fetch('/api/pipeline/status')
      .then(r => r.json())
      .then(status => {
        if (status.running) {
          updateLoadingStatus(status.step || 'Обработка данных...');
          setTimeout(() => checkPipelineStatus(), 1000);
        } else if (status.dataReady) {
          loadNews(); // Данные готовы, загружаем новости
        }
      })
      .catch(err => console.error('Ошибка проверки статуса:', err));
  }

  if (apply) {
    apply.addEventListener('click', () => {
      const filters = collectFiltersData(); 
      const welcome = document.querySelector('.welcome'); 
      if (welcome) welcome.style.display = 'none'; 
      showLoading(true);
      
      const params = new URLSearchParams();
      if (filters.categories && filters.categories.length) params.set('categories', filters.categories.join(','));
      const sentencesSelect = document.getElementById('sentences-select');
      const sentences = (sentencesSelect && sentencesSelect.value) ? sentencesSelect.value : '2';
      params.set('sentences', sentences);

      fetch('/api/refresh?' + params.toString(), { method: 'POST' })
        .then(r => { if (!r.ok) throw new Error('refresh failed'); return r.json(); })
        .then(() => fetch('/api/news?' + params.toString()))
        .then(r => { if (!r.ok) throw new Error('news fetch failed'); return r.json(); })
        .then(news => { showLoading(false); createMultipleNewsBlocks(Array.isArray(news) ? news : []); })
        .catch(err => { showLoading(false); console.error(err); alert('Ошибка при обновлении новостей'); });
    });
  }

  updateAuthButton();
  
  // Проверяем статус пайплайна при загрузке страницы
  setTimeout(() => checkPipelineStatus(), 500);
});