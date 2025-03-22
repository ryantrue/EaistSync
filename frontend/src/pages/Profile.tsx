// src/pages/Profile.tsx

import React, { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";

interface ProfileData {
    user_id: number;
    username: string;
    role: string;
    exp: number;
}

const Profile: React.FC = () => {
    const { accessToken } = useContext(AuthContext);
    const [profile, setProfile] = useState<ProfileData | null>(null);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        async function fetchProfile() {
            try {
                const response = await fetch("/api/profile", {
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": "Bearer " + accessToken
                    }
                });
                if (!response.ok) {
                    throw new Error("Ошибка получения профиля");
                }
                const data = await response.json();
                setProfile(data);
            } catch (err: any) {
                setError(err.message);
            }
        }

        if (accessToken) {
            fetchProfile();
        }
    }, [accessToken]);

    if (error) {
        return (
            <div className="container mt-5">
                <h2>Profile</h2>
                <div className="alert alert-danger">Ошибка: {error}</div>
            </div>
        );
    }

    if (!profile) {
        return (
            <div className="container mt-5">
                <h2>Profile</h2>
                <p>Загрузка профиля...</p>
            </div>
        );
    }

    return (
        <div className="container mt-5">
            <h2>Profile</h2>
            <p>
                <strong>ID:</strong> {profile.user_id}
            </p>
            <p>
                <strong>Username:</strong> {profile.username}
            </p>
            <p>
                <strong>Role:</strong> {profile.role}
            </p>
            <p>
                <strong>Expires at:</strong> {new Date(profile.exp * 1000).toLocaleString()}
            </p>
        </div>
    );
};

export default Profile;