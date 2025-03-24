import React from "react";
import { Row, Col, Form, Dropdown } from "react-bootstrap";
import SearchField from "./SearchField";
import { ContractFilters } from "../pages/Contracts";

export interface FiltersProps {
    filters: ContractFilters;
    suggestions: { [key in keyof Omit<ContractFilters, "state_Name">]: string[] };
    onFieldChange: (fieldKey: keyof ContractFilters, value: string | string[]) => void;
    onClearSuggestions: (fieldKey: keyof Omit<ContractFilters, "state_Name">) => void;
    selectSuggestion: (
        fieldKey: keyof Omit<ContractFilters, "state_Name">,
        suggestion: string
    ) => void;
    stateOptions: string[];
    statusCounts: Record<string, number>;
    totalCount: number;
}

/**
 * Компонент выбора статусов с пунктом "Все".
 * При загрузке все статусы выбраны (режим "Все" активен).
 * Если пользователь нажимает на "Все", когда оно активно – фильтр становится пустым (результат пуст).
 * Если нажимает, когда не активно – выбираются все статусы.
 * Выпадающий список закрывается при клике вне его (autoClose="outside").
 */
const MultiSelectStatus: React.FC<{
    selected: string[];
    options: string[];
    onChange: (selected: string[]) => void;
    statusCounts: Record<string, number>;
    totalCount: number;
}> = ({ selected, options, onChange, statusCounts, totalCount }) => {
    // Определяем, что "Все" активно, если в фильтре выбраны все опции.
    const allSelected = selected.length === options.length;
    // Для кнопки: если все выбраны – отображаем "Все", иначе – перечисляем выбранные статусы.
    const toggleText = allSelected ? "Все" : selected.join(", ");
    // Количество договоров для выбранных статусов: если выбраны все – общее число,
    // иначе сумма для выбранных статусов.
    const countDisplay = allSelected
        ? totalCount
        : selected.reduce((acc, cur) => acc + (statusCounts[cur] || 0), 0);

    // Обработчик клика для пункта "Все"
    const toggleAll = (e: React.MouseEvent<HTMLElement>) => {
        e.preventDefault();
        e.stopPropagation();
        // Если "Все" активно, то сбрасываем фильтр (пустой массив) – результат пустой,
        // иначе устанавливаем полный список опций.
        if (allSelected) {
            onChange([]);
        } else {
            onChange([...options]);
        }
    };

    // Обработчик клика для отдельной опции
    const toggleOption = (option: string, e: React.MouseEvent<HTMLElement>) => {
        e.preventDefault();
        e.stopPropagation();
        if (selected.includes(option)) {
            onChange(selected.filter((o) => o !== option));
        } else {
            onChange([...selected, option]);
        }
    };

    // Стиль для ограничения ширины кнопки переключения
    const toggleStyle: React.CSSProperties = {
        maxWidth: "200px",
        overflow: "hidden",
        textOverflow: "ellipsis",
        whiteSpace: "nowrap",
    };

    return (
        <Dropdown autoClose="outside">
            <Dropdown.Toggle variant="secondary" id="dropdown-status" style={toggleStyle}>
                {toggleText} ({countDisplay})
            </Dropdown.Toggle>
            <Dropdown.Menu style={{ minWidth: "300px", maxHeight: "300px", overflowY: "auto" }}>
                {/* Пункт "Все" */}
                <Dropdown.Item as="div" onClick={toggleAll}>
                    <Form.Check
                        type="checkbox"
                        label={`Все (${totalCount})`}
                        checked={allSelected}
                        readOnly
                    />
                </Dropdown.Item>
                {options.map((option, index) => (
                    <Dropdown.Item key={index} as="div" onClick={(e) => toggleOption(option, e)}>
                        <Form.Check
                            type="checkbox"
                            label={`${option} (${statusCounts[option] || 0})`}
                            checked={selected.includes(option)}
                            readOnly
                        />
                    </Dropdown.Item>
                ))}
            </Dropdown.Menu>
        </Dropdown>
    );
};

const filterFieldsConfig: Array<{
    key: keyof Omit<ContractFilters, "state_Name">;
    label: string;
    placeholder: string;
}> = [
    { key: "contractNumber", label: "Номер договора", placeholder: "Введите номер" },
    { key: "name", label: "Наименование", placeholder: "Введите наименование" },
    { key: "supplier_Name", label: "Контрагент", placeholder: "Введите контрагента" },
];

const Filters: React.FC<FiltersProps> = ({
                                             filters,
                                             suggestions,
                                             onFieldChange,
                                             onClearSuggestions,
                                             selectSuggestion,
                                             stateOptions,
                                             statusCounts,
                                             totalCount,
                                         }) => {
    return (
        <Form>
            <Row className="g-3">
                {filterFieldsConfig.map((field) => (
                    <Col md={3} key={field.key}>
                        <SearchField
                            fieldKey={field.key}
                            label={field.label}
                            placeholder={field.placeholder}
                            value={filters[field.key] as string}
                            suggestions={suggestions[field.key]}
                            onChange={onFieldChange}
                            onBlur={onClearSuggestions}
                            onSelectSuggestion={selectSuggestion}
                        />
                    </Col>
                ))}
                <Col md={3}>
                    <Form.Group controlId="filter-state">
                        <Form.Label>Статус:</Form.Label>
                        <MultiSelectStatus
                            selected={filters.state_Name}
                            options={stateOptions}
                            onChange={(selectedStatuses) =>
                                onFieldChange("state_Name", selectedStatuses)
                            }
                            statusCounts={statusCounts}
                            totalCount={totalCount}
                        />
                    </Form.Group>
                </Col>
            </Row>
        </Form>
    );
};

export default Filters;