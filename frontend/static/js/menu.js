/**
 * Асинхронно загружает HTML-разметку меню и вставляет её в контейнер.
 */
export async function loadMenu() {
    try {
        const response = await fetch("/static/components/menu.html");
        if (!response.ok) {
            throw new Error(`Ошибка загрузки меню: ${response.status}`);
        }
        const menuHTML = await response.text();
        const container = document.getElementById("menu-container");
        if (container) {
            container.innerHTML = menuHTML;
            initMenuEvents();
        }
    } catch (error) {
        console.error("Ошибка при загрузке меню:", error);
    }
}

/**
 * Инициализирует обработку событий меню.
 * Подсвечивает активный элемент в зависимости от текущего URL.
 */
function initMenuEvents() {
    const menuLinks = document.querySelectorAll(".navbar-nav .nav-link");
    const currentUrl = window.location.pathname;

    menuLinks.forEach(link => {
        const href = link.getAttribute("href");
        // Если URL совпадает или содержится в текущем пути, добавляем класс active
        if (href && currentUrl.includes(href)) {
            link.classList.add("active");
        }
        // При клике обновляем класс active для выбранного элемента
        link.addEventListener("click", () => {
            menuLinks.forEach(l => l.classList.remove("active"));
            link.classList.add("active");
        });
    });
}

// Загружаем меню, когда документ полностью загружен
document.addEventListener("DOMContentLoaded", loadMenu);