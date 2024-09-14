import { getHTTPHeaders, logout } from "./auth";
import {
	BundleModifiedResponse,
	CreateReleaseRequest,
	DeleteReleaseRequest,
	GenericResponse,
	ListAllBundlesResponse,
	ListAllReleasesResponse,
	OAuth2ConfigResponse,
	ReleaseModifiedResponse,
	SetReleaseActiveBundleRequest,
	UpdateReleaseRequest,
	UploadBundleRequest,
} from "./types";

const API_URL = "/api/v1";
const PUBLIC_API_URL = "/apipublic/v1";

const callApi = async (
	method: "POST" | "GET" | "PUT" | "DELETE",
	action: string,
	body: BodyInit | null = null,
	baseUrl: string = API_URL
) => {
	const response = await fetch(`${baseUrl}/${action}`, {
		method: method,
		headers: getHTTPHeaders(),
		body: body,
	});
	if (!response.ok) {
		if (response.status === 401) {
			logout()
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

export const getOAuth2Config = async (): Promise<OAuth2ConfigResponse> => {
	const response = await fetch(`${PUBLIC_API_URL}/oauth2.config`, {
		method: 'GET',
	});

	if (!response.ok) {
		throw new Error(`API Error ${response.status}: ${response.statusText}`);
	}

	return response.json();
};
