import React, { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { Container, Card, Table, Alert } from "react-bootstrap";
import { fetchContracts } from "../services/api";
import { logError, flattenContract } from "../utils/utils";
import Filters from "../components/Filters";
import fieldLabels from "../config/fieldLabels";

// Пример определения интерфейса контракта (расширьте по необходимости)
export interface Contract {
    id: number;
    contractNumber?: string;
    conclusionDate?: string;
    name?: string;
    supplier_Name?: string;
    state_Name?: string;
    [key: string]: any; // Для остальных полей
}

const allowedColumns: { [key: string]: string } = {
    contractNumber: fieldLabels.contractNumber,
    conclusionDate: fieldLabels.conclusionDate,
    name: fieldLabels.name,
    supplier_Name: fieldLabels.supplier_Name,
    state_Name: fieldLabels.state_Name,
};

const Contracts: React.FC = () => {
    const [contracts, setContracts] = useState<Contract[]>([]);
    const [filteredContracts, setFilteredContracts] = useState<Contract[]>([]);
    const [filters, setFilters] = useState({
        contractNumber: "",
        name: "",
        supplier_Name: "",
        state_Name: "",
    });
    const [suggestions, setSuggestions] = useState<{ [key: string]: string[] }>({
        contractNumber: [],
        name: [],
        supplier_Name: [],
    });
    const navigate = useNavigate();

    useEffect(() => {
        async function updateContracts() {
            try {
                const data = await fetchContracts();
                const flattened = data.map(flattenContract) as Contract[];
                setContracts(flattened);
                setFilteredContracts(flattened);
            } catch (error) {
                logError("Ошибка обновления контрактов", error as Error);
            }
        }
        updateContracts();
    }, []);

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

    const applyFilters = () => {
        const filtered = filterContracts();
        setFilteredContracts(filtered);
    };

    const getSuggestionsForField = useCallback(
        (fieldKey: string, value: string): string[] => {
            const filteredSuggestions = contracts.filter((contract) =>
                Object.entries(filters)
                    .filter(([key]) => key !== fieldKey)
                    .every(([key, filterValue]) => {
                        if (!filterValue) return true;
                        return contract[key] && contract[key].toLowerCase().includes(filterValue.toLowerCase());
                    }) && contract[fieldKey] && contract[fieldKey].toLowerCase().includes(value.toLowerCase())
            );
            return Array.from(new Set(filteredSuggestions.map((c) => c[fieldKey]))).sort((a, b) =>
                a.localeCompare(b)
            );
        },
        [contracts, filters]
    );

    const handleFieldChange = (fieldKey: string, value: string) => {
        setFilters((prev) => ({ ...prev, [fieldKey]: value }));
        if (["contractNumber", "name", "supplier_Name"].includes(fieldKey)) {
            setSuggestions((prev) => ({
                ...prev,
                [fieldKey]: getSuggestionsForField(fieldKey, value),
            }));
        }
    };

    const clearSuggestions = (fieldKey: string) => {
        setSuggestions((prev) => ({ ...prev, [fieldKey]: [] }));
    };

    const selectSuggestion = (fieldKey: string, suggestion: string) => {
        setFilters((prev) => ({ ...prev, [fieldKey]: suggestion }));
        setSuggestions((prev) => ({ ...prev, [fieldKey]: [] }));
    };

    const handleRowClick = (id: number) => {
        navigate(`/contract/${id}`);
    };

    const stateOptions = Array.from(
        new Set(contracts.map((contract) => contract.state_Name).filter((x): x is string => Boolean(x)))
    ).sort((a, b) => a.localeCompare(b));

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
                        onApplyFilters={applyFilters}
                    />
                </Card.Body>
            </Card>
            {filteredContracts.length === 0 ? (
                <Alert variant="info">Нет данных</Alert>
            ) : (
                <Table striped hover>
                    <thead>
                    <tr>
                        {Object.keys(allowedColumns).map((key) => (
                            <th key={key}>{allowedColumns[key]}</th>
                        ))}
                    </tr>
                    </thead>
                    <tbody>
                    {filteredContracts.map((contract) => (
                        <tr
                            key={contract.id}
                            onClick={() => handleRowClick(contract.id)}
                            style={{ cursor: "pointer" }}
                        >
                            {Object.keys(allowedColumns).map((key) => (
                                <td key={key}>
                                    {contract[key] !== undefined
                                        ? typeof contract[key] === "object" && contract[key] !== null
                                            ? JSON.stringify(contract[key])
                                            : contract[key]
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