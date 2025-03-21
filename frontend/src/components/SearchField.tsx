// src/components/SearchField.tsx
import React from "react";
import { Form } from "react-bootstrap";

interface SearchFieldProps<K extends string> {
    fieldKey: K;
    label: string;
    placeholder: string;
    value: string;
    suggestions: string[];
    onChange: (fieldKey: K, value: string) => void;
    onBlur: (fieldKey: K) => void;
    onSelectSuggestion: (fieldKey: K, suggestion: string) => void;
}

const SearchField = <K extends string>({
                                           fieldKey,
                                           label,
                                           placeholder,
                                           value,
                                           suggestions,
                                           onChange,
                                           onBlur,
                                           onSelectSuggestion,
                                       }: SearchFieldProps<K>) => {
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onChange(fieldKey, e.target.value);
    };

    return (
        <div className="position-relative">
            <Form.Group controlId={`filter-${fieldKey}`}>
                <Form.Label>{label}</Form.Label>
                <Form.Control
                    type="text"
                    placeholder={placeholder}
                    name={fieldKey}
                    value={value}
                    onChange={handleInputChange}
                    onBlur={() => setTimeout(() => onBlur(fieldKey), 150)}
                />
            </Form.Group>
            {suggestions && suggestions.length > 0 && (
                <div
                    className="dropdown-menu show"
                    style={{
                        width: "100%",
                        maxHeight: "200px",
                        overflowY: "auto",
                        position: "absolute",
                    }}
                >
                    {suggestions.map((suggestion, index) => (
                        <div
                            key={index}
                            className="dropdown-item"
                            style={{
                                whiteSpace: "nowrap",
                                overflow: "hidden",
                                textOverflow: "ellipsis",
                            }}
                            onMouseDown={() => onSelectSuggestion(fieldKey, suggestion)}
                        >
                            {suggestion}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default SearchField;