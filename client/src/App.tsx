import { Button, Layout, Menu, type MenuProps } from "antd";
import { useState } from "react";
import { Navigate, Outlet, useLocation, useNavigate } from "react-router-dom";
import type { MenuClickEventHandler } from "rc-menu/lib/interface";

const { Header, Content } = Layout;

type MenuItem = Required<MenuProps>["items"][number];
const items: MenuItem[] = [
	{
		key: "/app",
		label: "Releases",
	},
	{
		key: "/app/release/create",
		label: "Create Release",
	},
	{
		key: "/app/upload-bundle",
		label: "Upload Bundle",
	},
];

function App() {
	const token = localStorage.getItem("token");
	const navigate = useNavigate();
	const [currentPage, setCurrentPage] = useState(useLocation().pathname.replace(/\/$/, ""));
	const handleNavigate: MenuClickEventHandler = (e) => {
		setCurrentPage(e.key);
		navigate(e.key);
	};

	const handleLogout = () => {
		localStorage.removeItem("token");
		navigate("/");
	};

	return (
		<>
			{!token && <Navigate to="/login" />}
			<Layout>
				<Header style={{ display: "flex", alignItems: "center", gap: 48 }}>
					<h1 style={{ margin: 0, fontSize: 20, fontWeight: 600, color: "white" }}>Capgo UI</h1>
					<Menu
						theme="dark"
						selectedKeys={[currentPage]}
						onClick={handleNavigate}
						items={items}
						mode="horizontal"
						style={{ flex: 1, minWidth: 0 }}
					/>
					<Button type="text" style={{ color: "white" }} onClick={handleLogout}>
						Logout
					</Button>
				</Header>
				<Content style={{ padding: "0 24px", minHeight: "calc(100vh - 64px)" }}>
					<Outlet />
				</Content>
			</Layout>
		</>
	);
}

export default App;
