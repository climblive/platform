import { ApiClient, OrganizerCredentialsProvider } from "@climblive/lib";
import configData from "@climblive/lib/config.json";
import { navigate } from "svelte-routing";
import { SvelteDate } from "svelte/reactivity";
import { exchangeCode, refreshSession } from "./utils/cognito";

const checkTokensInterval = 60 * 1_000;
const minimumUsableTokenRemainingLifetime = 15 * 60 * 1_000;

export class Authenticator {
  private authenticated: boolean;
  private accessTokenExpiry: SvelteDate | undefined;
  private checkTokensIntervalTimer: number;

  constructor() {
    this.authenticated = $state(false);
    this.checkTokensIntervalTimer = 0;
  }

  public isAuthenticated = (): boolean => this.authenticated;

  public authenticate = async () => {
    const query = new URLSearchParams(location.search);
    const code = query.get("code");
    const state = query.get("state");

    if (code != null) {
      const { access_token, refresh_token } = await exchangeCode(code, state);

      ApiClient.getInstance().setCredentialsProvider(
        new OrganizerCredentialsProvider(access_token),
      );
      this.storeExpiryTime(access_token);

      localStorage.setItem("refresh_token", refresh_token);

      this.authenticated = true;

      navigate("./", { replace: true });

      return;
    }

    await this.refreshTokens();
  };

  private refreshTokens = async () => {
    if (
      this.accessTokenExpiry !== undefined &&
      this.accessTokenExpiry.getTime() - new SvelteDate().getTime() >=
        minimumUsableTokenRemainingLifetime
    ) {
      return;
    }

    try {
      const refreshToken = localStorage.getItem("refresh_token");

      if (refreshToken) {
        const { access_token } = await refreshSession(refreshToken);

        ApiClient.getInstance().setCredentialsProvider(
          new OrganizerCredentialsProvider(access_token),
        );
        this.storeExpiryTime(access_token);

        this.authenticated = true;
      }
    } catch {
      localStorage.removeItem("refresh_token");
      this.authenticated = false;
    }
  };

  private storeExpiryTime = (accessToken: string) => {
    const jwtPayload = JSON.parse(window.atob(accessToken.split(".")[1]));
    this.accessTokenExpiry = new SvelteDate(jwtPayload.exp * 1_000);
  };

  public startKeepAlive = () => {
    this.stopKeepAlive();

    this.refreshTokens();

    this.checkTokensIntervalTimer = setInterval(
      this.refreshTokens,
      checkTokensInterval,
    );
  };

  public stopKeepAlive = () => {
    if (this.checkTokensIntervalTimer) {
      clearInterval(this.checkTokensIntervalTimer);
      this.checkTokensIntervalTimer = 0;
    }
  };

  public redirectLogin = async () => {
    const codeVerifier = randomUrlSafeValue();
    const challenge = await challengeFromVerifier(codeVerifier);
    const state = randomUrlSafeValue();
    sessionStorage.setItem("oauth_code_verifier", codeVerifier);
    sessionStorage.setItem("oauth_state", state);

    const url = authorizationUrl("login", challenge, state);
    window.location.href = url;
  };

  public redirectSignup = async () => {
    const codeVerifier = randomUrlSafeValue();
    const challenge = await challengeFromVerifier(codeVerifier);
    const state = randomUrlSafeValue();
    sessionStorage.setItem("oauth_code_verifier", codeVerifier);
    sessionStorage.setItem("oauth_state", state);

    const url = authorizationUrl("signup", challenge, state);
    window.location.href = url;
  };

  public logout = () => {
    localStorage.removeItem("refresh_token");

    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/logout?client_id=${configData.COGNITO_CLIENT_ID}&logout_uri=${redirectUri}`;
    window.location.href = url;
  };
}

const randomUrlSafeValue = () => {
  const bytes = crypto.getRandomValues(new Uint8Array(32));
  return btoa(String.fromCharCode(...bytes))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
};

const authorizationUrl = (
  endpoint: "login" | "signup",
  challenge: string,
  state: string,
) => {
  const params = new URLSearchParams({
    response_type: "code",
    client_id: configData.COGNITO_CLIENT_ID,
    redirect_uri: window.location.origin + "/admin",
    code_challenge: challenge,
    code_challenge_method: "S256",
    state,
  });
  return `https://clmb.auth.eu-west-1.amazoncognito.com/${endpoint}?${params}`;
};

const challengeFromVerifier = async (verifier: string) => {
  const hash = await crypto.subtle.digest(
    "SHA-256",
    new TextEncoder().encode(verifier),
  );

  return btoa(String.fromCharCode(...new Uint8Array(hash)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
};
