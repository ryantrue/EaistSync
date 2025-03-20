/**
 * Получить значение параметра из URL.
 * @param {string} param - Название параметра.
 * @returns {string|null} Значение параметра или null, если параметр не найден.
 */
export function getQueryParam(param) {
    const params = new URLSearchParams(window.location.search);
    return params.get(param);
}

/**
 * Преобразует данные договора в плоскую структуру.
 * Если поле data представлено в виде строки, происходит попытка его парсинга.
 * @param {Object} item - Объект договора.
 * @returns {Object} Объединённый объект с данными и идентификатором.
 */
export function flattenContract(item) {
    let dataObj = {};
    if (item.data) {
        if (typeof item.data === "string") {
            try {
                dataObj = JSON.parse(item.data);
            } catch (e) {
                console.error("Ошибка парсинга data:", e);
            }
        } else if (typeof item.data === "object") {
            dataObj = item.data;
        }
    }
    return { ...dataObj, id: item.id };
}

/**
 * Централизованное логирование ошибок.
 * @param {string} message - Сообщение об ошибке.
 * @param {Error} error - Объект ошибки.
 */
export function logError(message, error) {
    console.error(message, error);
}