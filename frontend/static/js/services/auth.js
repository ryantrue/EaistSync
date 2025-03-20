// auth.js
// Этот скрипт отвечает за проверку авторизации и обработку форм авторизации/регистрации

document.addEventListener('DOMContentLoaded', () => {
    const authModalElement = document.getElementById('authModal');
    // Инициализируем модальное окно Bootstrap
    const authModal = new bootstrap.Modal(authModalElement, {
        backdrop: 'static',
        keyboard: false
    });

    // Проверка наличия токена (или другой логики аутентификации)
    function isAuthenticated() {
        // Пример: проверяем наличие authToken в localStorage
        return localStorage.getItem('authToken') !== null;
    }

    // Показываем модальное окно, если пользователь не авторизован
    if (!isAuthenticated()) {
        authModal.show();
    }

    // Обработка формы логина
    const loginForm = document.getElementById('loginForm');
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const username = document.getElementById('loginUsername').value;
        const password = document.getElementById('loginPassword').value;

        try {
            const response = await fetch('/login', { // Убедитесь, что путь соответствует вашему backend
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            if (!response.ok) {
                throw new Error('Ошибка авторизации');
            }

            const data = await response.json();
            console.log('Авторизация успешна:', data);

            // Сохраняем токен (если он возвращается) в localStorage
            if (data.token) {
                localStorage.setItem('authToken', data.token);
            }
            // Закрываем модальное окно
            authModal.hide();
        } catch (error) {
            console.error('Ошибка при авторизации:', error);
            alert('Ошибка авторизации. Проверьте введённые данные.');
        }
    });

    // Обработка формы регистрации
    const registerForm = document.getElementById('registerForm');
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const name = document.getElementById('registerName').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;
        const confirmPassword = document.getElementById('registerConfirmPassword').value;

        if (password !== confirmPassword) {
            alert('Пароли не совпадают!');
            return;
        }

        try {
            const response = await fetch('/register', { // Убедитесь, что путь соответствует вашему backend
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name, email, password })
            });

            if (!response.ok) {
                throw new Error('Ошибка регистрации');
            }

            const data = await response.json();
            console.log('Регистрация успешна:', data);

            // Если регистрация автоматически авторизует пользователя
            if (data.token) {
                localStorage.setItem('authToken', data.token);
            }
            authModal.hide();
        } catch (error) {
            console.error('Ошибка при регистрации:', error);
            alert('Ошибка регистрации. Попробуйте позже.');
        }
    });
});