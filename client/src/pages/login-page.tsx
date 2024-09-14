import { Button, Card, Divider, Flex, Input, Space } from "antd";
import { useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { isLoggedIn, oauthUserManager, setAPIToken, setOAuthToken } from "../client/auth";

export default function LoginPage() {
	const [apiKey, setApiKey] = useState("");
	const navigate = useNavigate();

	const handleLogin = () => {
		setAPIToken(apiKey);
		navigate(`${import.meta.env.BASE_URL}/app`);
	};

	const handleLoginWithOAuth2 = async () => {
		const user = await (await oauthUserManager).signinPopup({});
		setOAuthToken(user.access_token);
		navigate(`${import.meta.env.BASE_URL}/app`);
	};

	return (
		<>
			{isLoggedIn() && <Navigate to={`${import.meta.env.BASE_URL}/app`} />}
			<Flex justify="center" align="center" style={{ height: "100vh" }}>
				<Card style={{ minWidth: "350px" }}>
					<h1>Login</h1>

					<Space direction="vertical" style={{ width: "100%" }}>
						<p>Login with Management API Key</p>
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
					<Divider />
					<div>
						<Button type="primary" onClick={handleLoginWithOAuth2}>
							Login with OAuth 2
						</Button>
					</div>
				</Card>
			</Flex>
		</>
	);
}
