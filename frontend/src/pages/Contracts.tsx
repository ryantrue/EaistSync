// src/pages/Contracts.tsx
import React, { useState, useEffect, useCallback, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { Container, Card, Table, Alert, Spinner } from "react-bootstrap";
import { fetchContracts } from "../services/api";
import { logError, flattenContract } from "../utils/utils";
import Filters from "../components/Filters";
import fieldLabels from "../config/fieldLabels";

// Интерфейс для контракта (расширяйте по необходимости)
export interface Contract {
    id: number;
    contractNumber?: string;
    conclusionDate?: string;
    name?: string;
    supplier_Name?: string;
    state_Name?: string;
    [key: string]: any;
}

// Интерфейс для фильтров
export interface ContractFilters {
    contractNumber: string;
    name: string;
    supplier_Name: string;
    state_Name: string;
}

// Отображаемые колонки таблицы с подписями из конфигурации
const allowedColumns: { [key: string]: string } = {
    contractNumber: fieldLabels.contractNumber,
    conclusionDate: fieldLabels.conclusionDate,
    name: fieldLabels.name,
    supplier_Name: fieldLabels.supplier_Name,
    state_Name: fieldLabels.state_Name,
};

// Хук для дебаунсинга значения
function useDebounce<T>(value: T, delay: number): T {
    const [debouncedValue, setDebouncedValue] = useState<T>(value);

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedValue(value);
        }, delay);

        return () => clearTimeout(timer);
    }, [value, delay]);

    return debouncedValue;
}

const Contracts: React.FC = () => {
    const [contracts, setContracts] = useState<Contract[]>([]);
    const [filteredContracts, setFilteredContracts] = useState<Contract[]>([]);
    const [filters, setFilters] = useState<ContractFilters>({
        contractNumber: "",
        name: "",
        supplier_Name: "",
        state_Name: "",
    });
    // Подсказки только для текстовых полей (исключая state_Name)
    const [suggestions, setSuggestions] = useState<{
        [key in keyof Omit<ContractFilters, "state_Name">]: string[];
    }>({
        contractNumber: [],
        name: [],
        supplier_Name: [],
    });
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const navigate = useNavigate();

    // Загружаем данные при монтировании компонента
    useEffect(() => {
        const updateContracts = async () => {
            try {
                const data = await fetchContracts();
                const flattened = data.map(flattenContract) as Contract[];
                setContracts(flattened);
                setFilteredContracts(flattened);
            } catch (error) {
                logError("Ошибка обновления контрактов", error as Error);
            } finally {
                setIsLoading(false);
            }
        };
        updateContracts();
    }, []);

    // Функция фильтрации контрактов по заданным фильтрам
    const filterContracts = useCallback(() => {
        return contracts.filter((contract) =>
            Object.entries(filters).every(([key, filterValue]) => {
                if (!filterValue) return true;
                if (key === "state_Name") {
                    return (contract[key] || "").toLowerCase() === filterValue.toLowerCase();
                }
                return (contract[key] || "").toLowerCase().includes(filterValue.toLowerCase());
            })
        );
    }, [contracts, filters]);

    // Дебаунсим значения фильтров для оптимизации
    const debouncedFilters = useDebounce(filters, 300);

    // Применяем фильтрацию при изменении дебаунсённых фильтров
    useEffect(() => {
        setFilteredContracts(filterContracts());
    }, [debouncedFilters, filterContracts]);

    // Генерация подсказок для конкретного текстового поля с проверкой типа
    const getSuggestionsForField = useCallback(
        (fieldKey: keyof Omit<ContractFilters, "state_Name">, value: string): string[] => {
            const filteredSuggestions = contracts.filter((contract) => {
                const otherFiltersValid = Object.entries(filters)
                    .filter(([key]) => key !== fieldKey)
                    .every(([key, filterValue]) => {
                        if (!filterValue) return true;
                        if (key === "state_Name") {
                            return (contract[key] || "").toLowerCase() === filterValue.toLowerCase();
                        }
                        return (contract[key] || "").toLowerCase().includes(filterValue.toLowerCase());
                    });

                const fieldValue = contract[fieldKey];
                return (
                    otherFiltersValid &&
                    typeof fieldValue === "string" &&
                    fieldValue.toLowerCase().includes(value.toLowerCase())
                );
            });
            return Array.from(new Set(filteredSuggestions.map((c) => c[fieldKey] as string))).sort((a, b) =>
                a.localeCompare(b)
            );
        },
        [contracts, filters]
    );

    // Обработчик изменения значения поля фильтра
    const handleFieldChange = (fieldKey: keyof ContractFilters, value: string) => {
        setFilters((prev) => ({ ...prev, [fieldKey]: value }));
        if (fieldKey !== "state_Name") {
            setSuggestions((prev) => ({
                ...prev,
                [fieldKey]: getSuggestionsForField(fieldKey, value),
            }));
        }
    };

    // Очистка подсказок для конкретного поля
    const clearSuggestions = (fieldKey: keyof Omit<ContractFilters, "state_Name">) => {
        setSuggestions((prev) => ({ ...prev, [fieldKey]: [] }));
    };

    // Выбор подсказки для текстового поля
    const selectSuggestion = (fieldKey: keyof Omit<ContractFilters, "state_Name">, suggestion: string) => {
        setFilters((prev) => ({ ...prev, [fieldKey]: suggestion }));
        setSuggestions((prev) => ({ ...prev, [fieldKey]: [] }));
    };

    // Переход к деталям контракта по клику на строку таблицы
    const handleRowClick = (id: number) => {
        navigate(`/contract/${id}`);
    };

    // Вычисление уникальных статусов для фильтрации
    const stateOptions = useMemo(() => {
        return Array.from(
            new Set(contracts.map((contract) => contract.state_Name).filter((x): x is string => Boolean(x)))
        ).sort((a, b) => a.localeCompare(b));
    }, [contracts]);

    return (
        <Container className="my-4">
            <h1 className="mb-4">Контракты</h1>
            <Card className="mb-4">
                <Card.Body>
                    <Filters
                        filters={filters}
                        suggestions={suggestions}
                        onFieldChange={handleFieldChange}
                        onClearSuggestions={clearSuggestions}
                        selectSuggestion={selectSuggestion}
                        stateOptions={stateOptions}
                        // Можно оставить кнопку "Применить", если требуется ручной запуск фильтрации
                        onApplyFilters={() => setFilteredContracts(filterContracts())}
                    />
                </Card.Body>
            </Card>
            {isLoading ? (
                <div className="d-flex justify-content-center">
                    <Spinner animation="border" role="status">
                        <span className="visually-hidden">Загрузка...</span>
                    </Spinner>
                </div>
            ) : filteredContracts.length === 0 ? (
                <Alert variant="info">Нет данных</Alert>
            ) : (
                <Table striped hover>
                    <thead>
                    <tr>
                        {Object.keys(allowedColumns).map((colKey) => (
                            <th key={colKey}>{allowedColumns[colKey]}</th>
                        ))}
                    </tr>
                    </thead>
                    <tbody>
                    {filteredContracts.map((contract) => (
                        <tr key={contract.id} style={{ cursor: "pointer" }} onClick={() => handleRowClick(contract.id)}>
                            {Object.keys(allowedColumns).map((colKey) => (
                                <td key={colKey}>
                                    {contract[colKey] !== undefined
                                        ? typeof contract[colKey] === "object" && contract[colKey] !== null
                                            ? JSON.stringify(contract[colKey])
                                            : contract[colKey]
                                        : ""}
                                </td>
                            ))}
                        </tr>
                    ))}
                    </tbody>
                </Table>
            )}
        </Container>
    );
};

export default Contracts;