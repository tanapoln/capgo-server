import { Navigate } from "react-router-dom";

export default function LoginPage() {
    const token = localStorage.getItem("token");

    return (
        <>
        {token && <Navigate to="/app" />}
        <div>
            <h1>Login</h1>
        </div>
        </>
    )
}
