import { oauthUserManager } from "../client/auth";

export default function LoginOAuthCallbackPage() {
    (async () => {
        (await oauthUserManager).signinCallback(window.location.href)
    })()
	return <div>Login success</div>;
}
