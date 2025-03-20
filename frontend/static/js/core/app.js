import { initAuthModal } from "./authModal.js";
import { loadMenu } from "/static/js/menu.js";

document.addEventListener("DOMContentLoaded", async () => {
    // Загружаем меню
    await loadMenu();

    // Загружаем разметку модального окна авторизации
    try {
        const response = await fetch("/static/components/authModal.html");
        if (!response.ok) {
            throw new Error(`Ошибка загрузки authModal: ${response.status}`);
        }
        const authModalHTML = await response.text();
        const authContainer = document.getElementById("authModalContainer");
        if (authContainer) {
            authContainer.innerHTML = authModalHTML;
        }
    } catch (error) {
        console.error("Ошибка при загрузке authModal:", error);
    }

    // Инициализируем модальное окно авторизации
    initAuthModal();
});