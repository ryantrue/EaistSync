import React from "react";
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import Home from "./pages/Home";
import Contracts from "./pages/Contracts";
import States from "./pages/States";
import ContractDetail from "./pages/ContractDetail";
import { Container, Nav, Navbar } from "react-bootstrap";

const App: React.FC = () => {
    return (
        <Router>
            <Navbar bg="dark" variant="dark" expand="lg">
                <Container>
                    <Navbar.Brand as={Link} to="/">
                        EaistSync Dashboard
                    </Navbar.Brand>
                    <Navbar.Toggle aria-controls="basic-navbar-nav" />
                    <Navbar.Collapse id="basic-navbar-nav">
                        <Nav className="me-auto">
                            <Nav.Link as={Link} to="/">
                                Главная
                            </Nav.Link>
                            <Nav.Link as={Link} to="/contracts">
                                Контракты
                            </Nav.Link>
                            <Nav.Link as={Link} to="/states">
                                Состояния
                            </Nav.Link>
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>
            <main>
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/contracts" element={<Contracts />} />
                    <Route path="/states" element={<States />} />
                    <Route path="/contract/:id" element={<ContractDetail />} />
                </Routes>
            </main>
        </Router>
    );
};

export default App;