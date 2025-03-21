// src/utils/utils.ts

/**
 * Получить значение параметра из URL.
 * @param param - Название параметра.
 * @returns Значение параметра или null, если параметр не найден.
 */
export function getQueryParam(param: string): string | null {
    const params = new URLSearchParams(window.location.search);
    return params.get(param);
}

/**
 * Преобразует данные договора в плоскую структуру.
 * Если поле data представлено в виде строки, происходит попытка его парсинга.
 * @param item - Объект договора.
 * @returns Объединённый объект с данными и идентификатором.
 */
export function flattenContract(
    item: { data?: string | object; id: number; [key: string]: any }
): { id: number; [key: string]: any } {
    let dataObj: object = {};
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
 * @param message - Сообщение об ошибке.
 * @param error - Объект ошибки.
 */
export function logError(message: string, error: Error): void {
    console.error(message, error);
}