import React, { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { Container, Card, Table, Alert } from "react-bootstrap";
import { fetchContracts } from "../services/api.js";
import { logError, flattenContract } from "../utils/utils.js";
import Filters from "../components/Filters";
import allFieldLabels from "../config/fieldLabels";

// Для таблицы «Контракты» используем только выбранные поля:
const allowedColumns = {
    contractNumber: allFieldLabels.contractNumber,
    conclusionDate: allFieldLabels.conclusionDate,
    name: allFieldLabels.name,
    supplier_Name: allFieldLabels.supplier_Name,
    state_Name: allFieldLabels.state_Name,
};

function Contracts() {
    const [contracts, setContracts] = useState([]);
    const [filteredContracts, setFilteredContracts] = useState([]);
    const [filters, setFilters] = useState({
        contractNumber: "",
        name: "",
        supplier_Name: "",
        state_Name: "",
    });
    const [suggestions, setSuggestions] = useState({
        contractNumber: [],
        name: [],
        supplier_Name: [],
    });
    const navigate = useNavigate();

    useEffect(() => {
        async function updateContracts() {
            try {
                const data = await fetchContracts();
                const flattened = data.map(flattenContract);
                setContracts(flattened);
                setFilteredContracts(flattened);
            } catch (error) {
                logError("Ошибка обновления контрактов", error);
            }
        }
        updateContracts();
    }, []);

    // Универсальная фильтрация: для state_Name используется точное совпадение, для остальных – поиск по подстроке.
    const filterContracts = useCallback(() => {
        return contracts.filter((contract) =>
            Object.entries(filters).every(([key, filterValue]) => {
                if (!filterValue) return true;
                if (key === "state_Name") {
                    return contract[key] && contract[key].toLowerCase() === filterValue.toLowerCase();
                }
                return contract[key] && contract[key].toLowerCase().includes(filterValue.toLowerCase());
            })
        );
    }, [contracts, filters]);

    const applyFilters = () => {
        const filtered = filterContracts();
        setFilteredContracts(filtered);
    };

    // Универсальная функция формирования подсказок для текстового поля.
    const getSuggestionsForField = useCallback(
        (fieldKey, value) => {
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

    const handleFieldChange = (fieldKey, value) => {
        setFilters((prev) => ({ ...prev, [fieldKey]: value }));
        if (["contractNumber", "name", "supplier_Name"].includes(fieldKey)) {
            setSuggestions((prev) => ({
                ...prev,
                [fieldKey]: getSuggestionsForField(fieldKey, value),
            }));
        }
    };

    const clearSuggestions = (fieldKey) => {
        setSuggestions((prev) => ({ ...prev, [fieldKey]: [] }));
    };

    const selectSuggestion = (fieldKey, suggestion) => {
        setFilters((prev) => ({ ...prev, [fieldKey]: suggestion }));
        setSuggestions((prev) => ({ ...prev, [fieldKey]: [] }));
    };

    const handleRowClick = (id) => {
        navigate(`/contract/${id}`);
    };

    const stateOptions = Array.from(
        new Set(contracts.map((contract) => contract.state_Name).filter(Boolean))
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
}

export default Contracts;