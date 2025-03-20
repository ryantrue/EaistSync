import React from "react";
import ReactDOM from "react-dom/client";
import "./App.css"; // Ваши кастомные стили
import "bootstrap/dist/css/bootstrap.min.css"; // Bootstrap CSS
import App from "./App";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>
);