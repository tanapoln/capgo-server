import {
	BundleModifiedResponse,
	CreateReleaseRequest,
	DeleteReleaseRequest,
	GenericResponse,
	ListAllBundlesResponse,
	ListAllReleasesResponse,
	ReleaseModifiedResponse,
	SetReleaseActiveBundleRequest,
	UpdateReleaseRequest,
	UploadBundleRequest,
} from "./types";

const API_URL = "/api/v1";

const callApi = async (
	method: "POST" | "GET" | "PUT" | "DELETE",
	action: string,
	body: BodyInit | null = null
) => {
	const token = localStorage.getItem("token");
	const response = await fetch(`${API_URL}/${action}`, {
		method: method,
		headers: {
			"x-api-key": token ?? "",
		},
		body: body,
	});
	if (!response.ok) {
		if (response.status === 401) {
			localStorage.removeItem("token");
		}
		const errBody = await response.json();
		throw new Error(`API Error ${response.status}: ${errBody["error"] ?? errBody["message"] ?? ""}`);
	}
	return response;
};

export const listBundles = async (): Promise<ListAllBundlesResponse> => {
	const resp = await callApi("GET", "bundles.list");
	return resp.json();
};

export const uploadBundle = async (req: UploadBundleRequest): Promise<BundleModifiedResponse> => {
	const body = new FormData();
	body.append("bundle", req.bundle);
	body.append("app_id", req.app_id);
	body.append("version_name", req.version_name);
	body.append("description", req.description);

	const resp = await callApi("POST", "bundles.upload", body);
	return resp.json();
};

export const listReleases = async (): Promise<ListAllReleasesResponse> => {
	const resp = await callApi("GET", "releases.list");
	return resp.json();
};

export const createRelease = async (req: CreateReleaseRequest): Promise<ReleaseModifiedResponse> => {
	const resp = await callApi("POST", "releases.create", JSON.stringify(req));
	return resp.json();
};

export const updateRelease = async (req: UpdateReleaseRequest): Promise<ReleaseModifiedResponse> => {
	const resp = await callApi("POST", "releases.update", JSON.stringify(req));
	return resp.json();
};

export const deleteRelease = async (req: DeleteReleaseRequest): Promise<GenericResponse> => {
	const resp = await callApi("POST", "releases.delete", JSON.stringify(req));
	return resp.json();
};

export const setReleaseActiveBundle = async (
	req: SetReleaseActiveBundleRequest
): Promise<ReleaseModifiedResponse> => {
	const resp = await callApi("POST", "releases.set-active", JSON.stringify(req));
	return resp.json();
};
