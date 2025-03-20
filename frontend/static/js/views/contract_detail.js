import { fetchContracts } from "../services/api.js";
import { getQueryParam, flattenContract, logError } from "../utils/utils.js";

/**
 * Карта меток для отображения полей договора.
 * @type {Object<string, string>}
 */
const fieldLabels = {
    id: "ID",
    contractNumber: "Номер договора",
    conclusionDate: "Дата заключения",
    cost: "Стоимость",
    durationStartDate: "Начало действия",
    durationEndDate: "Окончание действия",
    name: "Наименование",
    supplier_Name: "Контрагент",
    state_Name: "Статус"
    // Дополнительные поля можно добавить здесь
};

/**
 * Получить детали договора по идентификатору.
 * Если отдельного эндпоинта нет, получаем все контракты и ищем нужный.
 * @param {string} contractId - Идентификатор договора.
 * @returns {Promise<Object|null>} - Обещание, возвращающее объект договора или null, если не найден.
 */
async function getContractDetail(contractId) {
    try {
        const data = await fetchContracts();
        const contracts = data.map(flattenContract);
        const contract = contracts.find(c => String(c.id) === String(contractId));
        return contract || null;
    } catch (error) {
        logError("Ошибка получения деталей договора:", error);
        return null;
    }
}

/**
 * Формирует HTML-таблицу с деталями договора.
 * Оформление с использованием классов Bootstrap.
 * @param {Object} data - Объект с данными договора.
 * @returns {string} - HTML-разметка таблицы.
 */
function renderContractDetail(data) {
    if (!data) return "<p>Данные не найдены</p>";

    let html = `
      <table class="table table-bordered table-striped">
        <thead>
          <tr>
            <th>Поле</th>
            <th>Значение</th>
          </tr>
        </thead>
        <tbody>
    `;
    Object.keys(data).forEach(key => {
        let value = data[key];
        if (typeof value === "object" && value !== null) {
            value = JSON.stringify(value);
        }
        const label = fieldLabels[key] || key;
        html += `<tr>
                   <td>${label}</td>
                   <td>${value !== undefined ? value : ""}</td>
                 </tr>`;
    });
    html += "</tbody></table>";
    return html;
}

/**
 * Инициализация страницы деталей договора.
 */
async function initContractDetailPage() {
    const contractId = getQueryParam("id");
    const detailContainer = document.getElementById("contract-detail");
    if (!contractId) {
        detailContainer.innerHTML = "<p>Не указан идентификатор договора.</p>";
        return;
    }
    const detail = await getContractDetail(contractId);
    detailContainer.innerHTML = renderContractDetail(detail);
}

document.addEventListener("DOMContentLoaded", initContractDetailPage);