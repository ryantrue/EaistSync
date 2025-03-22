// src/components/LazyLoader.tsx
import React, { Suspense } from "react";

interface LazyLoaderProps {
    fallback?: React.ReactNode;
    children: React.ReactNode;
}

const LazyLoader: React.FC<LazyLoaderProps> = ({ fallback = <div>Загрузка...</div>, children }) => {
    return <Suspense fallback={fallback}>{children}</Suspense>;
};

export default LazyLoader;