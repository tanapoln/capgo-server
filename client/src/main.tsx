import { App as AntApp } from "antd";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, Navigate, RouterProvider } from "react-router-dom";
import App from "./App.tsx";
import "./global.css";
import BundleUploadPage from "./pages/bundle-upload.tsx";
import LoginPage from "./pages/login-page.tsx";
import ReleaseCreatePage from "./pages/release-create.tsx";
import ReleaseUpdatePage from "./pages/release-update.tsx";
import ReleasesPage from "./pages/releases-page.tsx";

const router = createBrowserRouter([
	{
		path: "/",
		element: <RootPage />,
		children: [],
	},
	{
		path: "/login",
		element: <LoginPage />,
	},
	{
		path: "/app/",
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
]);

function RootPage() {
	const token = localStorage.getItem("token");
	return <>{token ? <Navigate to="/app" /> : <Navigate to="/login" />}</>;
}

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<AntApp>
			<RouterProvider router={router} />
		</AntApp>
	</StrictMode>
);
