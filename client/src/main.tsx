import { App as AntApp } from "antd";
import { StrictMode, useEffect } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, Outlet, RouterProvider, useLocation, useNavigate } from "react-router-dom";
import App from "./App.tsx";
import "./global.css";
import BundleUploadPage from "./pages/bundle-upload.tsx";
import LoginPage from "./pages/login-page.tsx";
import ReleaseCreatePage from "./pages/release-create.tsx";
import ReleaseUpdatePage from "./pages/release-update.tsx";
import ReleasesPage from "./pages/releases-page.tsx";

const router = createBrowserRouter([
	{
		path: import.meta.env.BASE_URL,
		element: <RootPage />,
		children: [
			{
				path: "login",
				element: <LoginPage />,
			},
			{
				path: "app",
				element: <App />,
				children: [
					{
						path: "",
						element: <ReleasesPage />,
					},
					{
						path: "release/create",
						element: <ReleaseCreatePage />,
					},
					{
						path: "release/:releaseId/update",
						element: <ReleaseUpdatePage />,
					},
					{
						path: "upload-bundle",
						element: <BundleUploadPage />,
					},
				],
			},
		],
	},
]);

// eslint-disable-next-line react-refresh/only-export-components
function RootPage() {
	const location = useLocation();
	const navigate = useNavigate();
	const token = localStorage.getItem("token");
	useEffect(() => {
		if (location.pathname.replace(/\/$/, "") === import.meta.env.BASE_URL) {
			if (token) {
				navigate("app");
			} else {
				navigate("login");
			}
		}
	}, [token, location.pathname, navigate])
	
	return <Outlet />;
}

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<AntApp>
			<RouterProvider router={router} />
		</AntApp>
	</StrictMode>
);
