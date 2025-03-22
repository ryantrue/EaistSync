// src/context/AuthContext.tsx

import React, { createContext, useState, useEffect, ReactNode } from 'react';
import { AuthResponse, login as apiLogin, logout as apiLogout, refreshToken as apiRefresh } from '../services/auth';

interface AuthContextProps {
    user: AuthResponse['user'] | null;
    accessToken: string | null;
    refreshToken: string | null;
    login: (username: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    refresh: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextProps>({
    user: null,
    accessToken: null,
    refreshToken: null,
    login: async () => {},
    logout: async () => {},
    refresh: async () => {}
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<AuthResponse['user'] | null>(null);
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const [refreshToken, setRefreshToken] = useState<string | null>(null);

    useEffect(() => {
        // Загружаем данные из localStorage (если есть)
        const storedUser = localStorage.getItem('user');
        const storedAccessToken = localStorage.getItem('accessToken');
        const storedRefreshToken = localStorage.getItem('refreshToken');
        if (storedUser && storedAccessToken && storedRefreshToken) {
            setUser(JSON.parse(storedUser));
            setAccessToken(storedAccessToken);
            setRefreshToken(storedRefreshToken);
        }
    }, []);

    const login = async (username: string, password: string) => {
        const authResponse = await apiLogin(username, password);
        setUser(authResponse.user);
        setAccessToken(authResponse.access_token);
        setRefreshToken(authResponse.refresh_token);
        localStorage.setItem('user', JSON.stringify(authResponse.user));
        localStorage.setItem('accessToken', authResponse.access_token);
        localStorage.setItem('refreshToken', authResponse.refresh_token);
    };

    const logout = async () => {
        if (refreshToken) {
            await apiLogout(refreshToken);
        }
        setUser(null);
        setAccessToken(null);
        setRefreshToken(null);
        localStorage.removeItem('user');
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
    };

    const refresh = async () => {
        if (refreshToken) {
            const response = await apiRefresh(refreshToken);
            setAccessToken(response.access_token);
            localStorage.setItem('accessToken', response.access_token);
        }
    };

    return (
        <AuthContext.Provider value={{ user, accessToken, refreshToken, login, logout, refresh }}>
            {children}
        </AuthContext.Provider>
    );
};