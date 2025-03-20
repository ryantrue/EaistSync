import React from "react";
import { Container } from "react-bootstrap";

function Home() {
    return (
        <Container className="my-4">
            <h1>Добро пожаловать в EaistSync Dashboard</h1>
            <p>Выберите раздел в меню для просмотра информации.</p>
        </Container>
    );
}

export default Home;