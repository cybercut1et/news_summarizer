// --- frontend/main.js (cleaned) ---
// UI helpers and server interaction
const loginBtn = document.getElementById('loginBtn');
const loginModal = document.getElementById('loginModal');
const loginForm = document.getElementById('loginForm');
const overlay = document.getElementById('overlay');
const usernameInput = document.getElementById('username');
const passwordInput = document.getElementById('password');

// Глобальные переменные для пагинации
let currentPage = 1;
let totalPages = 1;
let currentFilters = {};

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

  document.querySelectorAll('.categories-list button').forEach(b => b.addEventListener('click', () => {
    b.classList.toggle('active');
    renderFilteredNews();
  }));

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

  function createMultipleNewsBlocks(data) { 
    const container = document.querySelector('.news'); 
    if (!container) return; 
    container.innerHTML = '';
    
    // Если данные содержат пагинацию
    if (data.news && Array.isArray(data.news)) {
      data.news.forEach(createNewsBlock);
      updatePagination(data.pagination);
    } else if (Array.isArray(data)) {
      // Старый формат без пагинации
      data.forEach(createNewsBlock);
      hidePagination();
    }
  }

  function updatePagination(pagination) {
    if (!pagination) {
      hidePagination();
      return;
    }

    currentPage = pagination.currentPage;
    totalPages = pagination.totalPages;

    const paginationContainer = document.getElementById('pagination');
    const paginationInfo = document.getElementById('pagination-info-text');
    const prevBtn = document.getElementById('prev-page');
    const nextBtn = document.getElementById('next-page');
    const pageNumbers = document.getElementById('page-numbers');

    if (!paginationContainer) return;

    // Показываем пагинацию только если больше одной страницы
    if (totalPages <= 1) {
      hidePagination();
      return;
    }

    paginationContainer.style.display = 'block';

    // Обновляем информацию
    const start = (currentPage - 1) * pagination.limit + 1;
    const end = Math.min(currentPage * pagination.limit, pagination.totalItems);
    paginationInfo.textContent = `Показано ${start}-${end} из ${pagination.totalItems} новостей`;

    // Обновляем кнопки навигации
    prevBtn.disabled = !pagination.hasPrev;
    nextBtn.disabled = !pagination.hasNext;

    // Генерируем номера страниц
    pageNumbers.innerHTML = '';
    const maxVisiblePages = 5;
    let startPage = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2));
    let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);

    if (endPage - startPage < maxVisiblePages - 1) {
      startPage = Math.max(1, endPage - maxVisiblePages + 1);
    }

    for (let i = startPage; i <= endPage; i++) {
      const pageBtn = document.createElement('button');
      pageBtn.textContent = i;
      pageBtn.type = 'button';
      pageBtn.className = i === currentPage ? 'active' : '';
      pageBtn.addEventListener('click', () => loadNewsPage(i));
      pageNumbers.appendChild(pageBtn);
    }
  }

  function hidePagination() {
    const paginationContainer = document.getElementById('pagination');
    if (paginationContainer) {
      paginationContainer.style.display = 'none';
    }
  }

  function loadNewsPage(page = 1) {
    currentPage = page;
    loadNews(currentPage);
  }

  function collectFiltersData() { 
    const urls = []; 
    if (containerUrl) containerUrl.querySelectorAll('input[type="url"]').forEach(i => { if (i.value.trim()) urls.push(i.value.trim()); }); 
    // Маппим выбранные русские категории к английским
    const catsRu = Array.from(document.querySelectorAll('.categories-list button.active')).map(b => b.textContent.trim().toLowerCase());
    const catsEn = catsRu.map(cat => categoryMap[cat]).filter(Boolean);
    const periodEl = document.getElementById('tailmetr'); 
    const period = periodEl ? parseInt(periodEl.value) : 24; 
    return { urls, categories: catsEn, period }; 
  }

  function loadNews(page = 1) {
    const filters = collectFiltersData(); 
    currentFilters = filters; // Сохраняем текущие фильтры
    const welcome = document.querySelector('.welcome'); 
    if (welcome) welcome.style.display = 'none'; 
    showLoading(true);
    
    const params = new URLSearchParams();
    if (filters.categories && filters.categories.length) params.set('categories', filters.categories.join(','));
    const sentencesSelect = document.getElementById('sentences-select');
    const sentences = (sentencesSelect && sentencesSelect.value) ? sentencesSelect.value : '2';
    params.set('sentences', sentences);
    params.set('page', page.toString());
    params.set('limit', '10'); // Устанавливаем лимит в 10 новостей на страницу

    fetch('/api/news?' + params.toString())
      .then(r => { if (!r.ok) throw new Error('news fetch failed'); return r.json(); })
      .then(data => { 
        // Проверяем, загружаются ли данные
        if (data.status === 'loading') {
          // Показываем статус пайплайна и проверяем снова
          updateLoadingStatus(data.pipeline_status || 'Обработка данных...');
          setTimeout(() => loadNews(page), 2000); // Проверяем каждые 2 секунды
          return;
        }
        
        showLoading(false); 
        if ((data.news && Array.isArray(data.news) && data.news.length > 0) || (Array.isArray(data) && data.length > 0)) {
          createMultipleNewsBlocks(data);
        } else {
          const w = document.querySelector('.welcome');
          if (w) {
            w.innerHTML = 'Новостей не найдено. Данные обрабатываются...';
            w.style.display = 'block';
          }
          hidePagination();
        }
      })
      .catch(err => { showLoading(false); hidePagination(); console.error(err); alert('Ошибка при загрузке новостей'); });
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
          loadNews(1); // Данные готовы, загружаем первую страницу новостей
        }
      })
      .catch(err => console.error('Ошибка проверки статуса:', err));
  }

  // Обработчики для кнопок пагинации
  const prevBtn = document.getElementById('prev-page');
  const nextBtn = document.getElementById('next-page');
  
  if (prevBtn) {
    prevBtn.addEventListener('click', () => {
      if (currentPage > 1) {
        loadNewsPage(currentPage - 1);
      }
    });
  }
  
  if (nextBtn) {
    nextBtn.addEventListener('click', () => {
      if (currentPage < totalPages) {
        loadNewsPage(currentPage + 1);
      }
    });
  }

  if (apply) {
    apply.addEventListener('click', () => {
      currentPage = 1; // Сбрасываем на первую страницу при применении фильтров
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
        .then(() => {
          // После обновления загружаем первую страницу
          const newsParams = new URLSearchParams();
          if (filters.categories && filters.categories.length) newsParams.set('categories', filters.categories.join(','));
          newsParams.set('sentences', sentences);
          newsParams.set('page', '1');
          newsParams.set('limit', '10');
          
          return fetch('/api/news?' + newsParams.toString());
        })
        .then(r => { if (!r.ok) throw new Error('news fetch failed'); return r.json(); })
        .then(data => { 
          showLoading(false); 
          createMultipleNewsBlocks(data);
        })
        .catch(err => { showLoading(false); hidePagination(); console.error(err); alert('Ошибка при обновлении новостей'); });
    });
  }

  let allNewsData = [];

  function fetchFilteredData() {
    return fetch('/ml/filtered_data.json')
      .then(r => r.json())
      .then(data => {
        allNewsData = [];
        data.forEach(channel => {
          channel.messages.forEach(msg => {
            allNewsData.push({
              channel_name: channel.channel_name,
              ...msg
            });
          });
        });
        renderFilteredNews();
      })
      .catch(err => {
        console.error('Ошибка загрузки filtered_data.json:', err);
        document.querySelector('.news').innerHTML = '<div style="color:red">Ошибка загрузки новостей</div>';
      });
  }


  // Маппинг русских названий категорий к английским значениям из filtered_data.json
  const categoryMap = {
    'глянец': 'gloss',
    'здоровье': 'health',
    'климат': 'climate',
    'конфликты': 'conflicts',
    'культура': 'culture',
    'экономика': 'economy',
    'наука': 'science',
    'общество': 'society',
    'политика': 'politics',
    'спорт': 'sports',
    'путешествия': 'travel'
  };

  function renderFilteredNews() {
    const activeCatsRu = Array.from(document.querySelectorAll('.categories-list button.active')).map(b => b.textContent.trim().toLowerCase());
    // Маппим выбранные русские категории к английским
    const activeCatsEn = activeCatsRu.map(cat => categoryMap[cat]).filter(Boolean);
    let filtered = allNewsData;
    if (activeCatsEn.length > 0) {
      filtered = filtered.filter(n => n.category && activeCatsEn.includes(n.category.toLowerCase()));
    }
    // Если ни одна категория не выбрана, показываем все новости
    const container = document.querySelector('.news');
    container.innerHTML = '';
    if (filtered.length === 0) {
      container.innerHTML = '<div style="color:gray">Нет новостей по выбранным категориям</div>';
      return;
    }
    filtered.forEach(item => {
      const block = document.createElement('div');
      block.className = 'class-block-news';
      block.innerHTML = `${item.category ? `<div class="cat">${item.category}</div>` : ''}<div class="text">${item.text}</div>${item.link ? `<a class="link" href="${item.link}" target="_blank">Перейти в источник</a>` : ''}`;
      container.appendChild(block);
    });
  }

  updateAuthButton();
  
  // Проверяем статус пайплайна при загрузке страницы
  setTimeout(() => checkPipelineStatus(), 500);
  fetchFilteredData();
});