// src/services/api.js

/**
 * Универсальная функция для получения JSON по указанному URL.
 * @param {string} url - URL для запроса.
 * @returns {Promise<any>} - Обещание, которое возвращает распарсенный JSON.
 */
export async function fetchJSON(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Ошибка HTTP: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error(`Ошибка запроса к ${url}:`, error);
        throw error;
    }
}

/**
 * Получить список контрактов.
 * @returns {Promise<any>} - Обещание, возвращающее данные контрактов.
 */
export async function fetchContracts() {
    return await fetchJSON("/api/contracts");
}

/**
 * Получить список состояний.
 * @returns {Promise<any>} - Обещание, возвращающее данные состояний.
 */
export async function fetchStates() {
    return await fetchJSON("/api/states");
}

/**
 * Получить список лотов.
 * @returns {Promise<any>} - Обещание, возвращающее данные лотов.
 */
export async function fetchLots() {
    return await fetchJSON("/api/lots");
}