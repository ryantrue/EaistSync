// src/components/Filters.tsx
import React from "react";
import { Row, Col, Form, Button } from "react-bootstrap";
import SearchField from "./SearchField";
import { ContractFilters } from "../pages/Contracts";

export interface FiltersProps {
    filters: ContractFilters;
    suggestions: { [key in keyof Omit<ContractFilters, "state_Name">]: string[] };
    onFieldChange: (fieldKey: keyof ContractFilters, value: string) => void;
    onClearSuggestions: (fieldKey: keyof Omit<ContractFilters, "state_Name">) => void;
    selectSuggestion: (fieldKey: keyof Omit<ContractFilters, "state_Name">, suggestion: string) => void;
    stateOptions: string[];
    onApplyFilters: () => void;
}

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
                                             onApplyFilters,
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
                            value={filters[field.key]}
                            suggestions={suggestions[field.key]}
                            onChange={onFieldChange}
                            onBlur={onClearSuggestions}
                            onSelectSuggestion={selectSuggestion}
                        />
                    </Col>
                ))}
                <Col md={2}>
                    <Form.Group controlId="filter-state">
                        <Form.Label>Статус:</Form.Label>
                        <Form.Select
                            value={filters.state_Name}
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