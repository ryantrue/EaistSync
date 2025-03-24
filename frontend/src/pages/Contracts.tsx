import React, { useState, useEffect, useCallback, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { Container, Card, Table, Alert, Spinner } from "react-bootstrap";
import { fetchContracts } from "../services/api";
import { logError, flattenContract } from "../utils/utils";
import Filters from "../components/Filters";
import fieldLabels from "../config/fieldLabels";

// Интерфейс для контракта
export interface Contract {
    id: number;
    contractNumber?: string;
    conclusionDate?: string;
    name?: string;
    supplier_Name?: string;
    state_Name?: string;
    [key: string]: any;
}

// Интерфейс для фильтров: поле state_Name – массив строк
export interface ContractFilters {
    contractNumber: string;
    name: string;
    supplier_Name: string;
    state_Name: string[];
}

// Отображаемые колонки таблицы
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
        const timer = setTimeout(() => setDebouncedValue(value), delay);
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
        state_Name: [],
    });
    const [suggestions, setSuggestions] = useState<{
        [key in keyof Omit<ContractFilters, "state_Name">]: string[];
    }>({
        contractNumber: [],
        name: [],
        supplier_Name: [],
    });
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [isStatusInitialized, setIsStatusInitialized] = useState(false);
    const navigate = useNavigate();

    // Загрузка договоров
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

    // Вычисление уникальных статусов
    const stateOptions = useMemo(() => {
        return Array.from(
            new Set(
                contracts
                    .map((contract) => contract.state_Name)
                    .filter((x): x is string => Boolean(x))
            )
        ).sort((a, b) => a.localeCompare(b));
    }, [contracts]);

    // При первой загрузке, если фильтр по статусу не установлен, устанавливаем его как полный список (режим "Все")
    useEffect(() => {
        if (!isStatusInitialized && stateOptions.length > 0) {
            setFilters((prev) => ({ ...prev, state_Name: stateOptions }));
            setIsStatusInitialized(true);
        }
    }, [stateOptions, isStatusInitialized]);

    /**
     * Функция фильтрации:
     * Для полей, кроме state_Name, если значение фильтра пустое – контракт проходит фильтрацию.
     * Для state_Name, если массив пустой – контракт НЕ проходит фильтрацию (результат пуст).
     */
    const filterContracts = useCallback(() => {
        return contracts.filter((contract) =>
            Object.entries(filters).every(([key, filterValue]) => {
                if (key === "state_Name") {
                    if (Array.isArray(filterValue) && filterValue.length === 0) return false;
                    return (filterValue as string[]).some(
                        (status) => (contract[key] || "").toLowerCase() === status.toLowerCase()
                    );
                } else {
                    if (!filterValue || (Array.isArray(filterValue) && filterValue.length === 0))
                        return true;
                    return (contract[key] || "").toLowerCase().includes(
                        (filterValue as string).toLowerCase()
                    );
                }
            })
        );
    }, [contracts, filters]);

    const debouncedFilters = useDebounce(filters, 300);
    useEffect(() => {
        setFilteredContracts(filterContracts());
    }, [debouncedFilters, filterContracts]);

    const getSuggestionsForField = useCallback(
        (fieldKey: keyof Omit<ContractFilters, "state_Name">, value: string): string[] => {
            const filteredSuggestions = contracts.filter((contract) => {
                const otherFiltersValid = Object.entries(filters)
                    .filter(([key]) => key !== fieldKey)
                    .every(([key, filterValue]) => {
                        if (!filterValue || (Array.isArray(filterValue) && filterValue.length === 0))
                            return true;
                        if (key === "state_Name") {
                            return (filterValue as string[]).some(
                                (status) => (contract[key] || "").toLowerCase() === status.toLowerCase()
                            );
                        }
                        return (contract[key] || "").toLowerCase().includes(
                            (filterValue as string).toLowerCase()
                        );
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

    const handleFieldChange = (fieldKey: keyof ContractFilters, value: string | string[]) => {
        setFilters((prev) => ({ ...prev, [fieldKey]: value }));
        if (fieldKey !== "state_Name") {
            setSuggestions((prev) => ({
                ...prev,
                [fieldKey]: getSuggestionsForField(fieldKey, value as string),
            }));
        }
    };

    const clearSuggestions = (fieldKey: keyof Omit<ContractFilters, "state_Name">) => {
        setSuggestions((prev) => ({ ...prev, [fieldKey]: [] }));
    };

    const selectSuggestion = (fieldKey: keyof Omit<ContractFilters, "state_Name">, suggestion: string) => {
        setFilters((prev) => ({ ...prev, [fieldKey]: suggestion }));
        setSuggestions((prev) => ({ ...prev, [fieldKey]: [] }));
    };

    const handleRowClick = (id: number) => {
        navigate(`/contract/${id}`);
    };

    const statusCounts = useMemo(() => {
        const counts: Record<string, number> = {};
        contracts.forEach((contract) => {
            const status = contract.state_Name;
            if (status) {
                counts[status] = (counts[status] || 0) + 1;
            }
        });
        return counts;
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
                        statusCounts={statusCounts}
                        totalCount={contracts.length}
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
                        <tr
                            key={contract.id}
                            style={{ cursor: "pointer" }}
                            onClick={() => handleRowClick(contract.id)}
                        >
                            {Object.keys(allowedColumns).map((colKey) => (
                                <td key={colKey}>
                                    {contract[colKey] !== undefined
                                        ? typeof contract[colKey] === "object" &&
                                        contract[colKey] !== null
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