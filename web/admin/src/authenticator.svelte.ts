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

    if (code != null) {
      const { access_token, refresh_token } = await exchangeCode(code);

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
    const verifier = generateRandomString();
    const challenge = await challenge_from_verifier(verifier);
    sessionStorage.setItem("code_verifier", verifier);

    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/login?response_type=code&client_id=${configData.COGNITO_CLIENT_ID}&redirect_uri=${redirectUri}&code_challenge=${challenge}&code_challenge_method=S256`;
    window.location.href = url;
  };

  public redirectSignup = () => {
    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/signup?response_type=code&client_id=${configData.COGNITO_CLIENT_ID}&redirect_uri=${redirectUri}`;
    window.location.href = url;
  };

  public logout = () => {
    localStorage.removeItem("refresh_token");

    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/logout?client_id=${configData.COGNITO_CLIENT_ID}&logout_uri=${redirectUri}`;
    window.location.href = url;
  };
}

function dec2hex(dec: number) {
  return ("0" + dec.toString(16)).substr(-2);
}

function generateRandomString() {
  const array = new Uint32Array(56 / 2);
  window.crypto.getRandomValues(array);
  return Array.from(array, dec2hex).join("");
}

function sha256(plain: string): Promise<ArrayBuffer> {
  const encoder = new TextEncoder();
  const data = encoder.encode(plain);
  return window.crypto.subtle.digest("SHA-256", data);
}

function base64urlencode(a: ArrayBuffer) {
  let str = "";
  const bytes = new Uint8Array(a);
  const len = bytes.byteLength;

  for (let i = 0; i < len; i++) {
    str += String.fromCharCode(bytes[i]);
  }

  return btoa(str).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

async function challenge_from_verifier(v: string) {
  const hashed = await sha256(v);
  const base64encoded = base64urlencode(hashed);

  return base64encoded;
}
