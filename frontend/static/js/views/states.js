import { fetchStates } from "../services/api.js";
import { logError } from "../utils/utils.js";

/**
 * Формирует HTML-таблицу для данных состояний.
 * Оформление с использованием классов Bootstrap.
 * @param {Array<Object>} data - Массив объектов состояний.
 * @returns {string} HTML-разметка таблицы.
 */
function renderTable(data) {
    if (!data || data.length === 0) {
        return "<p>Нет данных</p>";
    }

    let html = `
      <table class="table table-striped table-hover">
        <thead>
          <tr>`;
    // Заголовки таблицы получаются из ключей первого объекта
    Object.keys(data[0]).forEach(key => {
        html += `<th>${key}</th>`;
    });
    html += `  </tr>
        </thead>
        <tbody>`;
    data.forEach(item => {
        html += "<tr>";
        for (const key in item) {
            let value = item[key];
            if (typeof value === "object" && value !== null) {
                value = JSON.stringify(value);
            }
            html += `<td>${value !== undefined ? value : ""}</td>`;
        }
        html += "</tr>";
    });
    html += "</tbody></table>";
    return html;
}

/**
 * Обновляет данные состояний и выводит их в контейнер с id="states".
 */
async function updateStates() {
    try {
        const states = await fetchStates();
        const statesDiv = document.getElementById("states");
        if (statesDiv) {
            statesDiv.innerHTML = renderTable(states);
        }
    } catch (error) {
        logError("Ошибка обновления состояний:", error);
    }
}

/**
 * Инициализация страницы состояний.
 * Обновление данных происходит сразу и затем каждые 30 секунд.
 */
function initStatesPage() {
    updateStates();
    setInterval(updateStates, 30000);
}

document.addEventListener("DOMContentLoaded", initStatesPage);