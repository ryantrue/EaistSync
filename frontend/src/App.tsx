// src/App.tsx
import React, { useContext, lazy } from "react";
import { BrowserRouter as Router, Routes, Route, Link, Navigate } from "react-router-dom";
import LazyLoader from "./components/LazyLoader";
import { Container, Nav, Navbar } from "react-bootstrap";
import { AuthContext, AuthProvider } from "./context/AuthContext";

const Home = lazy(() => import("./pages/Home"));
const Contracts = lazy(() => import("./pages/Contracts"));
const States = lazy(() => import("./pages/States"));
const ContractDetail = lazy(() => import("./pages/ContractDetail"));
const Login = lazy(() => import("./pages/Login"));
const Register = lazy(() => import("./pages/Register"));
const Profile = lazy(() => import("./pages/Profile"));

// Компонент для защиты маршрутов
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { user } = useContext(AuthContext);
    return user ? <>{children}</> : <Navigate to="/" replace />;
};

const App: React.FC = () => {
    const { user, logout } = useContext(AuthContext);

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
                            {user && (
                                <>
                                    <Nav.Link as={Link} to="/contracts">
                                        Контракты
                                    </Nav.Link>
                                    <Nav.Link as={Link} to="/states">
                                        Состояния
                                    </Nav.Link>
                                    <Nav.Link as={Link} to="/profile">
                                        Профиль
                                    </Nav.Link>
                                </>
                            )}
                        </Nav>
                        <Nav>
                            {user ? (
                                <>
                                    <Nav.Link disabled>Привет, {user.username}</Nav.Link>
                                    <Nav.Link onClick={logout}>Logout</Nav.Link>
                                </>
                            ) : (
                                <Nav.Link as={Link} to="/register">
                                    Register
                                </Nav.Link>
                            )}
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>
            <main>
                <LazyLoader>
                    <Routes>
                        <Route path="/" element={user ? <Home /> : <Login />} />
                        <Route path="/register" element={<Register />} />
                        <Route
                            path="/contracts"
                            element={
                                <ProtectedRoute>
                                    <Contracts />
                                </ProtectedRoute>
                            }
                        />
                        <Route
                            path="/states"
                            element={
                                <ProtectedRoute>
                                    <States />
                                </ProtectedRoute>
                            }
                        />
                        <Route
                            path="/contract/:id"
                            element={
                                <ProtectedRoute>
                                    <ContractDetail />
                                </ProtectedRoute>
                            }
                        />
                        <Route
                            path="/profile"
                            element={
                                <ProtectedRoute>
                                    <Profile />
                                </ProtectedRoute>
                            }
                        />
                        <Route path="*" element={<Navigate to="/" replace />} />
                    </Routes>
                </LazyLoader>
            </main>
        </Router>
    );
};

const AppWithProvider: React.FC = () => (
    <AuthProvider>
        <App />
    </AuthProvider>
);

export default AppWithProvider;