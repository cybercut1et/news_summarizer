const loginBtn = document.getElementById('loginBtn');
const loginModal = document.getElementById('loginModal');
const loginForm = document.getElementById('loginForm');
const overlay = document.getElementById('overlay');
const usernameInput = document.getElementById('username');
const passwordInput = document.getElementById('password');
const cancel=document.getElementById('cancel');
function openModal() {
    loginModal.style.display = 'flex';
    overlay.style.display = 'block';
    document.body.style.overflow = 'hidden'; // Блокируем прокрутку страницы
}

// Функция закрытия модального окна
function closeModal() {
    loginModal.style.display = 'none';
    overlay.style.display = 'none';
    document.body.style.overflow = ''; // Разблокируем прокрутку
    loginForm.reset(); // Очищаем форму
}

loginBtn.addEventListener('click', openModal);
cancel.addEventListener('click', closeModal);

document.addEventListener('keydown', function(event) {
    // Проверяем, что нажата клавиша Enter
    if (event.key === 'Esc') {
       event.preventDefault();
        // Можно добавить проверку, что модальное окно действительно открыто
        // Например, если у вас есть переменная-флаг или проверка CSS класса
        closeModal();
    }
});


function updateAuthButton() {
    const loginButton = document.getElementById('loginBtn');
    const username = localStorage.getItem('username');
    
    if (username) {
        // После входа - показываем картинку
        loginButton.innerHTML = '<img src="../frontend/icons/free-icon-font-user-3917711.png" style="width:24px; height: auto;">';
        loginButton.className = 'avatar-btn';
    } else {
        // До входа - показываем текст "Вход"
        loginButton.innerHTML = 'Вход';
        loginButton.className = 'login-btn';
    }
}


// Закрытие по кнопке "Отправить"
loginForm.addEventListener('submit', function(event) {
    event.preventDefault(); // Предотвращаем перезагрузку страницы
    // Получаем значения полей
    const username = usernameInput.value;
    const password = passwordInput.value;
    
    // Проверяем, что поля не пустые
    // if (username.trim() === '' || password.trim() === '') {
    //     alert('Пожалуйста, заполните все поля');
    //     return;
    // }
// Сохраняем данные
    localStorage.setItem('username', username);
    localStorage.setItem('password', password);
    
    // Можно отправить данные на сервер (пример с fetch)
    // sendToServer(username, password);
    
    console.log('Данные для входа:', { username, password });
    
    // Закрываем модальное окно
    updateAuthButton(); // Обновляем кнопку входа
    closeModal();
    
});

overlay.addEventListener('click', closeModal);

// Закрытие по Escape
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        closeModal();
    }
});

let windowCount = 1;
document.getElementById('plus').addEventListener('click', function() {
    // Создаем новое окошечко
    const newWindow = document.createElement('div');
    newWindow.className = 'url-box';
    
    // Добавляем содержимое в окошечко
    newWindow.innerHTML = '<input type="url" placeholder="ссылка на ресурс">';
    
    // Добавляем окошечко в контейнер
    document.getElementById('container-url').appendChild(newWindow);

    windowCount++;
    
    // Выводим в консоль текущее количество (для отладки)
    // console.log('Создано окошечек: ' + (windowCount - 1));
});

document.getElementById('minus').addEventListener('click', function() {
    const windows = document.querySelectorAll('.url-box');
    
    if (windows.length > 1) {
        const lastWindow = windows[windows.length - 1];
         lastWindow.classList.add('removing');
        // Анимация удаления
        // lastWindow.style.animation = 'slideOut 0.3s ease-out';
        
        setTimeout(() => {
            lastWindow.remove();
            windowCount = windows.length - 1; // обновляем счетчик
        },300);
    } 
});



const buttons = document.querySelectorAll('.categories-list button');


buttons.forEach(button => {
        button.addEventListener('click', function() {
            this.classList.toggle('active');
        });
    });

    const range = document.querySelector('.range-input')
range.addEventListener('input', handleInputRange)

function handleInputRange() {
  event.target.parentNode.parentNode.style.setProperty(
    '--value',
    event.target.value
  )
}
function handleInputRange() {
  event.target.parentNode.parentNode.style.setProperty(
    '--value',
    event.target.value
  )
  // изменение значения тега `<output>`
  event.target.nextElementSibling.value = event.target.value;
}
 


//эксперименты с отправкой на сервак!!!!!!

document.addEventListener('DOMContentLoaded', function() {
    const applyButton = document.getElementById('apply');
    const rangeInput = document.getElementById('tailmetr');
    const rangeOutput = document.getElementById('output');
    const plusButton = document.getElementById('plus');
    const minusButton = document.getElementById('minus');
    const containerUrl = document.getElementById('container-url');

    // Обновление значения range
    rangeInput.addEventListener('input', function() {
        rangeOutput.textContent = this.value;
    });

    

    // Обработка нажатия на кнопку "Применить"  КНОПКА!!!!1
   applyButton.addEventListener('click', function() {
    const filtersData = collectFiltersData();
    console.log('Данные для отправки:', filtersData);
    // Удаляем класс у приветственного окна
    const welcomeElement = document.querySelector('.welcome'); 
    if (welcomeElement) {
          welcomeElement.style.display = 'none';
    }

    showLoading(true);
//ТЕСТОВЫЙ ВАРИАНТ (пока нет сервера)
    testServerDelay()
        .then(() => {
            // 4. Скрываем загрузку
            showLoading(false);


    createTestNewsBlock()});   //createMultipleNewsBlocks();


// РЕАЛЬНЫЙ ВАРИАНТ ЗАГРУЗКИ (когда будет сервер)
    /*
    fetchFromServer(filtersData)
        .then(result => {
            // 4. Скрываем загрузку
            showLoading(false);
            
            // 5. Показываем реальные новости с сервера
            const newsArray = result.news || result.articles || result.data || [];
            createMultipleNewsBlocks(newsArray);
        })
        .catch(error => {
            // Скрываем загрузку даже при ошибке
            showLoading(false);
            console.error('Ошибка:', error);
        });
    */




    // Отправка на реальный сервер
    fetch('/api/get-news', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(filtersData)
    })
    .then(response => response.json())
    .then(result => {
       
        createMultipleNewsBlocks(result.news);   //result.news--предполагаемое название массива новостей с сервера
    })
    .catch(error => {
        console.error('Ошибка:', error);
    });
});

    // Сбор всех данных фильтров
    function collectFiltersData() {
        return {
            urls: getUrls(),
            categories: getActiveCategories(),
            period: getPeriodValue(),
            // Добавь другие параметры если нужно
        };
    }

    // Получение всех URL из полей ввода
    function getUrls() {
        const urlInputs = containerUrl.querySelectorAll('input[type="url"]');
        const urls = [];
        
        urlInputs.forEach(input => {
            if (input.value.trim()) {
                urls.push(input.value.trim());
            }
        });
        
        return urls;
    }

    // Получение активных категорий
    function getActiveCategories() {
        const activeCategories = [];
        const categoryButtons = document.querySelectorAll('.categories-list button');
        
        categoryButtons.forEach(button => {
            // Проверяем активна ли кнопка (можно добавить класс active при клике)
            if (button.classList.contains('active')) {
                activeCategories.push(button.textContent.trim());
            }
        });
        
        return activeCategories;
    }

    // Получение значения периода
    function getPeriodValue() {
        return parseInt(rangeInput.value);
    }

    // Функция для отображения данных (для тестирования)
    function displayRequestData(requestData) {
        // Можно вывести в console или показать на странице
        // console.log('Данные для отправки:', JSON.parse(requestData.body));
        console.log('Данные для отправки:', JSON.parse(JSON.stringify(filtersData)));
        // Или создать блок для отображения
        const outputDiv = document.getElementById('output') || createOutputDiv();
        outputDiv.innerHTML = '<h4>Данные для отправки:</h4><pre>' + 
                             JSON.stringify(JSON.parse(requestData.body), null, 2) + '</pre>';
    }

    // Функция для создания одной новости с данными с сервера
function createNewsBlock(newsData) {
    const newsBlock = document.createElement('div');
    newsBlock.className = 'class-block-news';
    
    newsBlock.innerHTML = `
        <div class="text">
            ${newsData.text || 'текст новости'}       
        </div>
        <div class="date">${newsData.date || '01.01.2001'} 
            <div class="time">${newsData.time || '22:48'}</div>
        </div>
        ${newsData.category ? `<div class="cat">${newsData.category}</div>` : ''}
        ${newsData.source ? `<a class="link" href="${newsData.source}" target="_blank">Перейти в источник</a>` : ''}
    `;
    
    const newsContainer = document.querySelector('.news');
    if (newsContainer) {
        newsContainer.appendChild(newsBlock);
    }
    
    return newsBlock;
}

// Функция для создания нескольких новостей
function createMultipleNewsBlocks(newsArray) {
    // Очищаем старые новости
    const newsContainer = document.querySelector('.news');
    if (newsContainer) {
        newsContainer.innerHTML = '';
    }
    
    // Создаем блоки для каждой новости
    newsArray.forEach((newsItem) => {
        createNewsBlock(newsItem);
    });
}



//ТЕСТОВАЯ ФУНКЦИЯ КОТОРАЯ ПРОСТО СОЗДАЁТ НПС БЛОК ДЛЯ ДЕМОНСТРАЦИИ
function createTestNewsBlock() {
    // Создаем основной блок новости
    const newsBlock = document.createElement('div');
    newsBlock.className = 'class-block-news';
    
    // Заполняем содержимое как в твоем примере
    newsBlock.innerHTML = `
    <div class="cat">
                    Категория
                   </div>
        <div class="text">
            Текст новости  
            </div>
            <a class="link" href="#" target="_blank">Перейти в источник</a>
        <div class="date">01.01.2001 
            <div class="time">22:48</div>
        </div>
    `;
    
    // Находим контейнер для новостей и добавляем созданный блок
    const newsContainer = document.querySelector('.news');
    if (newsContainer) {
        newsContainer.appendChild(newsBlock);
    } else {
        console.error('Контейнер для новостей не найден!');
    }
    
    return newsBlock;
}



function sendToServer(requestData) {
        fetch(requestData.url, {
            method: requestData.method,
            headers: requestData.headers,
            body: requestData.body
        })
        .then(response => response.json())
        .then(data => {
            console.log('Ответ от сервера:', data);
            // Обработка ответа от сервера
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
    }
});

///ДЛЯ ЭКРАНА ЗАГРУЗКИ 


// Показывает/скрывает загрузку
function showLoading(show) {
    const welcomeElement = document.querySelector('.welcome');
    
    if (!welcomeElement) return;
    
    if (show) {
        // Превращаем welcome в экран загрузки
        welcomeElement.innerHTML = `
            <div class="loading">
                <img src="../frontend/icons/free-icon-old-scroll-17472284.png" alt="Загрузка" style="width:110px; height: auto;" class="loading-image">
                <div class="loading-text">Идет поиск новостей...</div>
            </div>
        `;
        welcomeElement.style.display = 'block';
    } else {
        // Скрываем welcome (или можно вернуть исходное содержимое)
        welcomeElement.style.display = 'none';
    }
}



// function createLoadingElement() {
//     const loadingDiv = document.createElement('div');
//     loadingDiv.id = 'loading';
//     loadingDiv.className = 'loading';
//     loadingDiv.innerHTML = `
//         <div class="loading">
//             <img src="../frontend/icons/free-icon-old-scroll-17472284.png" style="width:60px; height: auto;" alt="Загрузка" class="loading-image">
//             <div class="loading-text">Идет поиск новостей...</div>
//         </div>
//     `;
//     document.body.appendChild(loadingDiv);
//     return loadingDiv;
// }


// ТЕСТОВАЯ функция - только имитирует задержку сервера
function testServerDelay() {
    return new Promise((resolve) => {
        // Имитация задержки сервера (2 секунды)
        setTimeout(() => {
            resolve();
        }, 2000);
    });
}


// РЕАЛЬНАЯ функция - только запрос к серверу
function fetchFromServer(filtersData) {
    return fetch('/api/get-news', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(filtersData)
    })
    .then(response => response.json());
}