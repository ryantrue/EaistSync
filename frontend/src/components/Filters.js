import React from "react";
import { Row, Col, Form, Button } from "react-bootstrap";

function Filters({
                     filters,
                     suggestions,
                     onFieldChange,      // Универсальный обработчик изменения значения поля
                     onClearSuggestions, // Функция для очистки подсказок для поля
                     selectSuggestion,   // Функция выбора подсказки
                     stateOptions,       // Массив вариантов для поля «Статус»
                     onApplyFilters,     // Функция применения фильтров
                 }) {
    // Универсальный обработчик для текстовых полей
    const getInputHandlers = (fieldKey) => ({
        value: filters[fieldKey] || "",
        onChange: (e) => onFieldChange(fieldKey, e.target.value),
        onBlur: () => setTimeout(() => onClearSuggestions(fieldKey), 150),
    });

    // Конфигурация текстовых полей
    const filterFields = [
        { key: "contractNumber", label: "Номер договора", placeholder: "Введите номер" },
        { key: "name", label: "Наименование", placeholder: "Введите наименование" },
        { key: "supplier_Name", label: "Контрагент", placeholder: "Введите контрагента" },
    ];

    return (
        <Form>
            <Row className="g-3">
                {filterFields.map(({ key, label, placeholder }) => (
                    <Col md={3} className="position-relative" key={key}>
                        <Form.Group controlId={`filter-${key}`}>
                            <Form.Label>{label}</Form.Label>
                            <Form.Control
                                type="text"
                                placeholder={placeholder}
                                {...getInputHandlers(key)}
                            />
                        </Form.Group>
                        {suggestions[key] && suggestions[key].length > 0 && (
                            <div
                                className="dropdown-menu show"
                                style={{
                                    width: "100%",
                                    maxHeight: "200px",
                                    overflowY: "auto",
                                    position: "absolute",
                                }}
                            >
                                {suggestions[key].map((item, index) => (
                                    <div
                                        key={index}
                                        className="dropdown-item"
                                        style={{
                                            whiteSpace: "nowrap",
                                            overflow: "hidden",
                                            textOverflow: "ellipsis",
                                        }}
                                        onMouseDown={() => selectSuggestion(key, item)}
                                    >
                                        {item}
                                    </div>
                                ))}
                            </div>
                        )}
                    </Col>
                ))}

                {/* Фильтр по статусу */}
                <Col md={2}>
                    <Form.Group controlId="filter-state">
                        <Form.Label>Статус:</Form.Label>
                        <Form.Select
                            value={filters.state_Name || ""}
                            onChange={(e) => onFieldChange("state_Name", e.target.value)}
                        >
                            <option value="">-- Выберите статус --</option>
                            {stateOptions.map((status, index) => (
                                <option key={index} value={status}>
                                    {status}
                                </option>
                            ))}
                        </Form.Select>
                    </Form.Group>
                </Col>

                <Col md={1} className="d-flex align-items-end">
                    <Button variant="primary" onClick={onApplyFilters}>
                        Применить
                    </Button>
                </Col>
            </Row>
        </Form>
    );
}

export default Filters;