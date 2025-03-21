// src/pages/ContractDetail.tsx
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Container, Table, Spinner, Alert } from "react-bootstrap";
import { fetchContracts } from "../services/api";
import { flattenContract, logError } from "../utils/utils";
import fieldLabels from "../config/fieldLabels";

interface Contract {
    id: number;
    [key: string]: any;
}

const ContractDetail: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const [contract, setContract] = useState<Contract | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        async function loadContract() {
            try {
                const data = await fetchContracts();
                const contracts: Contract[] = data.map(flattenContract);
                const found = contracts.find((c: Contract) => String(c.id) === String(id));
                setContract(found || null);
            } catch (err) {
                // Приводим err к типу Error
                logError("Ошибка получения деталей договора:", err as Error);
                setError("Ошибка загрузки данных");
            } finally {
                setLoading(false);
            }
        }
        loadContract();
    }, [id]);

    if (loading) {
        return (
            <Container className="my-4 text-center">
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Загрузка...</span>
                </Spinner>
            </Container>
        );
    }

    if (error) {
        return (
            <Container className="my-4">
                <Alert variant="danger">{error}</Alert>
            </Container>
        );
    }

    if (!contract) {
        return (
            <Container className="my-4">
                <Alert variant="warning">Данные не найдены</Alert>
            </Container>
        );
    }

    return (
        <Container className="my-4">
            <h1 className="mb-4">Детали договора</h1>
            <Table striped bordered hover>
                <thead>
                <tr>
                    <th>Поле</th>
                    <th>Значение</th>
                </tr>
                </thead>
                <tbody>
                {Object.keys(contract).map((key) => {
                    let value = contract[key];
                    if (typeof value === "object" && value !== null) {
                        value = JSON.stringify(value);
                    }
                    const label = fieldLabels[key] || key;
                    return (
                        <tr key={key}>
                            <td>{label}</td>
                            <td>{value !== undefined ? value : ""}</td>
                        </tr>
                    );
                })}
                </tbody>
            </Table>
        </Container>
    );
};

export default ContractDetail;