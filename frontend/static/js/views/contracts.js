import { fetchContracts } from "../services/api.js";
import { logError, flattenContract } from "../utils/utils.js";

/**
 * Определение столбцов для отображения и их меток.
 * @type {Object<string, string>}
 */
const allowedColumns = {
    contractNumber: "Номер договора",
    conclusionDate: "Дата заключения",
    name: "Наименование",
    supplier_Name: "Контрагент",
    state_Name: "Статус"
};

let allContracts = [];

/**
 * Рендерит HTML-таблицу контрактов.
 * @param {Array<Object>} data - Массив контрактов.
 * @returns {string} - HTML-разметка таблицы.
 */
function renderTable(data) {
    if (!data || data.length === 0) return "<p>Нет данных</p>";

    let html = `<table class="table table-striped table-hover">
                  <thead>
                    <tr>`;
    for (let key in allowedColumns) {
        html += `<th>${allowedColumns[key]}</th>`;
    }
    html += `   </tr>
                </thead>
                <tbody>`;
    data.forEach(item => {
        html += `<tr data-id="${item.id}" class="clickable-row">`;
        for (let key in allowedColumns) {
            let value = item[key];
            if (typeof value === "object" && value !== null) {
                value = JSON.stringify(value);
            }
            html += `<td>${value !== undefined ? value : ""}</td>`;
        }
        html += "</tr>";
    });
    html += `  </tbody>
              </table>`;
    return html;
}

/**
 * Собирает значения фильтров из полей расширенного поиска.
 * @returns {Object} - Значения фильтров.
 */
function getAdvancedFilters() {
    return {
        contractNumber: document.getElementById("filter-contractNumber").value.trim().toLowerCase(),
        name: document.getElementById("filter-name").value.trim().toLowerCase(),
        supplier_Name: document.getElementById("filter-supplier").value.trim().toLowerCase(),
        state_Name: document.getElementById("filter-state").value.trim().toLowerCase()
    };
}

/**
 * Фильтрует контракты на основе введённых значений.
 * @returns {Array<Object>} - Отфильтрованный список контрактов.
 */
function filterContracts() {
    const filters = getAdvancedFilters();
    return allContracts.filter(contract => {
        let ok = true;
        if (filters.contractNumber &&
            (!contract.contractNumber || contract.contractNumber.toLowerCase().indexOf(filters.contractNumber) === -1)) {
            ok = false;
        }
        if (filters.name &&
            (!contract.name || contract.name.toLowerCase().indexOf(filters.name) === -1)) {
            ok = false;
        }
        if (filters.supplier_Name &&
            (!contract.supplier_Name || contract.supplier_Name.toLowerCase().indexOf(filters.supplier_Name) === -1)) {
            ok = false;
        }
        if (filters.state_Name &&
            (!contract.state_Name || contract.state_Name.toLowerCase() !== filters.state_Name)) {
            ok = false;
        }
        return ok;
    });
}

/**
 * Обновляет таблицу контрактов согласно фильтрам.
 */
function applyFilters() {
    const filtered = filterContracts();
    const container = document.getElementById("contracts");
    container.innerHTML = renderTable(filtered);
    addRowClickHandlers();
}

/**
 * Добавляет обработчики клика для строк таблицы.
 */
function addRowClickHandlers() {
    document.querySelectorAll("tr.clickable-row").forEach(row => {
        row.addEventListener("click", () => {
            const id = row.getAttribute("data-id");
            // Переход на страницу деталей договора с передачей id в query-параметрах
            window.location.href = `../../pages/contract_detail.html?id=${id}`;
        });
    });
}

/**
 * Скрывает все открытые выпадающие списки.
 */
function hideAllDropdowns() {
    document.querySelectorAll(".dropdown").forEach(dropdown => {
        dropdown.style.display = "none";
    });
}

/**
 * Обновляет выпадающий список для указанного текстового поля.
 * @param {string} fieldKey - "contractNumber", "name" или "supplier_Name".
 */
function updateDropdownForField(fieldKey) {
    hideAllDropdowns();

    let inputId, dropdownId;
    switch (fieldKey) {
        case "contractNumber":
            inputId = "filter-contractNumber";
            dropdownId = "dropdown-contractNumber";
            break;
        case "name":
            inputId = "filter-name";
            dropdownId = "dropdown-name";
            break;
        case "supplier_Name":
            inputId = "filter-supplier";
            dropdownId = "dropdown-supplier";
            break;
        default:
            return;
    }
    const inputElement = document.getElementById(inputId);
    const dropdownElement = document.getElementById(dropdownId);
    const text = inputElement.value.trim().toLowerCase();

    // Задаем ширину выпадающего списка равной ширине поля ввода
    const inputWidth = inputElement.offsetWidth;
    dropdownElement.style.width = inputWidth + "px";

    // Получаем фильтры без текущего поля
    const filters = getAdvancedFilters();
    filters[fieldKey] = "";
    const suggestions = allContracts.filter(contract => {
        let ok = true;
        for (const key in filters) {
            if (key === fieldKey) continue;
            if (filters[key] && (!contract[key] || contract[key].toLowerCase().indexOf(filters[key]) === -1)) {
                ok = false;
                break;
            }
        }
        return ok && contract[fieldKey] && contract[fieldKey].toLowerCase().includes(text);
    });

    // Формируем уникальный список вариантов для данного поля
    const uniqueValues = Array.from(new Set(suggestions.map(contract => contract[fieldKey]).filter(Boolean)))
        .sort((a, b) => a.localeCompare(b));

    // Генерируем HTML для элементов выпадающего списка
    let html = uniqueValues.map(value =>
        `<div class="dropdown-item" style="
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            max-width: 100%;
        ">${value}</div>`
    ).join("");

    dropdownElement.innerHTML = html;
    dropdownElement.style.maxHeight = "200px";
    dropdownElement.style.overflowY = "auto";
    dropdownElement.style.display = html ? "block" : "none";

    // Обработка выбора элемента из списка
    dropdownElement.querySelectorAll(".dropdown-item").forEach(item => {
        item.addEventListener("mousedown", () => {
            inputElement.value = item.textContent;
            dropdownElement.innerHTML = "";
            dropdownElement.style.display = "none";
        });
    });
}

/**
 * Заполняет select для статуса уникальными значениями из списка контрактов.
 */
function populateStatusSelect() {
    const selectElement = document.getElementById("filter-state");
    selectElement.innerHTML = `<option value="">-- Выберите статус --</option>`;
    const uniqueStatuses = Array.from(new Set(allContracts
        .map(contract => contract.state_Name)
        .filter(Boolean)))
        .sort((a, b) => a.localeCompare(b));
    uniqueStatuses.forEach(value => {
        const option = document.createElement("option");
        option.value = value;
        option.textContent = value;
        selectElement.appendChild(option);
    });
}

/**
 * Загружает данные контрактов и инициализирует страницу.
 */
async function updateContracts() {
    try {
        const data = await fetchContracts();
        allContracts = data.map(flattenContract);
        populateStatusSelect();
        const container = document.getElementById("contracts");
        container.innerHTML = renderTable(allContracts);
        addRowClickHandlers();
    } catch (error) {
        logError("Ошибка обновления контрактов", error);
    }
}

document.addEventListener("DOMContentLoaded", () => {
    // Обработчики для полей фильтра с динамическими выпадающими списками
    const textFields = [
        { key: "contractNumber", inputId: "filter-contractNumber", dropdownId: "dropdown-contractNumber" },
        { key: "name", inputId: "filter-name", dropdownId: "dropdown-name" },
        { key: "supplier_Name", inputId: "filter-supplier", dropdownId: "dropdown-supplier" }
    ];

    textFields.forEach(field => {
        const input = document.getElementById(field.inputId);
        if (input) {
            input.addEventListener("focus", () => updateDropdownForField(field.key));
            input.addEventListener("input", () => updateDropdownForField(field.key));
            input.addEventListener("blur", () => {
                setTimeout(() => {
                    const dropdown = document.getElementById(field.dropdownId);
                    if (dropdown) {
                        dropdown.style.display = "none";
                    }
                }, 150);
            });
        }
    });

    // Глобальный обработчик клика для закрытия выпадающих списков при клике вне области фильтра
    document.addEventListener("click", e => {
        if (!e.target.closest(".filter-field") && !e.target.closest(".dropdown")) {
            hideAllDropdowns();
        }
    });

    // Обработчик кнопки «Применить фильтры»
    const applyFiltersBtn = document.getElementById("apply-filters");
    if (applyFiltersBtn) {
        applyFiltersBtn.addEventListener("click", applyFilters);
    }

    // Инициализируем загрузку данных контрактов
    updateContracts();
});