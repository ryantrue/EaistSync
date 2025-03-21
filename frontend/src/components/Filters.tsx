// src/components/Filters.tsx
import React from "react";
import { Row, Col, Form, Button } from "react-bootstrap";

interface FiltersProps {
    filters: { [key: string]: string };
    suggestions: { [key: string]: string[] };
    onFieldChange: (fieldKey: string, value: string) => void;
    onClearSuggestions: (fieldKey: string) => void;
    selectSuggestion: (fieldKey: string, suggestion: string) => void;
    stateOptions: string[];
    onApplyFilters: () => void;
}

const Filters: React.FC<FiltersProps> = ({
                                             filters,
                                             suggestions,
                                             onFieldChange,
                                             onClearSuggestions,
                                             selectSuggestion,
                                             stateOptions,
                                             onApplyFilters,
                                         }) => {
    // Универсальный обработчик для текстовых полей
    const getInputHandlers = (fieldKey: string) => ({
        value: filters[fieldKey] || "",
        onChange: (e: React.ChangeEvent<HTMLInputElement>) => onFieldChange(fieldKey, e.target.value),
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
                            <Form.Control type="text" placeholder={placeholder} {...getInputHandlers(key)} />
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

                <Col md={2}>
                    <Form.Group controlId="filter-state">
                        <Form.Label>Статус:</Form.Label>
                        <Form.Select
                            value={filters.state_Name || ""}
                            onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
                                onFieldChange("state_Name", e.target.value)
                            }
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
};

export default Filters;