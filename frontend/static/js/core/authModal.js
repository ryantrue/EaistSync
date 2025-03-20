// authModal.js
export function initAuthModal() {
    const authModalElement = document.getElementById('authModal');
    if (!authModalElement) {
        console.error('Элемент модального окна авторизации (authModal) не найден в DOM.');
        return;
    }

    // Инициализируем модальное окно Bootstrap
    const authModal = new bootstrap.Modal(authModalElement, {
        backdrop: 'static',
        keyboard: false
    });

    // Функция проверки авторизации: здесь можно проверить наличие токена или иной признак авторизации
    function isAuthenticated() {
        return localStorage.getItem('authToken') !== null;
    }

    if (!isAuthenticated()) {
        authModal.show();
    }

    // Возвращаем объект с методами для управления модальным окном при необходимости
    return {
        show: () => authModal.show(),
        hide: () => authModal.hide()
    };
}