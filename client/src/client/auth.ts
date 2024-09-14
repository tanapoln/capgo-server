import { UserManager } from "oidc-client-ts";
import { getOAuth2Config } from "./api";

export const oauthUserManager = (async () => {
	const config = await getOAuth2Config();
	return new UserManager({
		authority: config.issuer,
		client_id: config.client_id,
		redirect_uri: redirectUrl(),
		response_type: "code",
		scope: "openid",
	});
})()

export const isLoggedIn = () => {
    return localStorage.getItem("token") || localStorage.getItem("oauth_token")
}

export const logout = () => {
    localStorage.removeItem("token")
    localStorage.removeItem("oauth_token")
}

export const setAPIToken = (token: string) => {
    localStorage.setItem("token", token)
}

export const setOAuthToken = (token: string) => {
    localStorage.setItem("oauth_token", token)
}

export const getHTTPHeaders = () => {
	const headers = new Headers()
	
	const token = localStorage.getItem("token")
	if (token) {
		headers.set("x-api-key", token)
	}

	const oauthToken = localStorage.getItem("oauth_token")
	if (oauthToken) {
		headers.set("Authorization", `Bearer ${oauthToken}`)
	}

    return headers
}

function redirectUrl() {
	const url = new URL(window.location.href)
	url.pathname = "/ui/login/oauth-callback"
	return url.toString()
}
