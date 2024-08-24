import { Button, Card, Flex, Input, Space } from "antd";
import { useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";

export default function LoginPage() {
	const token = localStorage.getItem("token");
	const [apiKey, setApiKey] = useState("");
	const navigate = useNavigate();

	const handleLogin = () => {
		localStorage.setItem("token", apiKey);
		navigate(`${import.meta.env.BASE_URL}/app`);
	};

	return (
		<>
			{token && <Navigate to={`${import.meta.env.BASE_URL}/app`} />}
			<Flex justify="center" align="center" style={{ height: "100vh" }}>
				<Card style={{ minWidth: "350px" }}>
					<h1>Login</h1>
					<Space direction="vertical" style={{ width: "100%" }}>
						<Input
							placeholder="Capgo server API Key"
							value={apiKey}
							onChange={(e) => setApiKey(e.target.value)}
							onPressEnter={handleLogin}
						/>
						<Button type="primary" onClick={handleLogin}>
							Login
						</Button>
					</Space>
				</Card>
			</Flex>
		</>
	);
}
