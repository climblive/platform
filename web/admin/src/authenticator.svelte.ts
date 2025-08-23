import { ApiClient, OrganizerCredentialsProvider } from "@climblive/lib";
import configData from "@climblive/lib/config.json";
import { navigate } from "svelte-routing";
import { SvelteDate } from "svelte/reactivity";
import { exchangeCode, refreshSession } from "./utils/cognito";

const checkTokensInterval = 60 * 1_000;
const minimumUsableTokenRemainingLifetime = 15 * 60 * 1_000;

export class Authenticator {
  private authenticated = $state(false);
  private accessTokenExpiry: SvelteDate | undefined;
  private checkTokensIntervalTimer: number = 0;

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
      this.accessTokenExpiry.getTime() - new Date().getTime() >=
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
    this.accessTokenExpiry = new Date(jwtPayload.exp * 1_000);
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

  public redirectLogin = () => {
    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/login?response_type=code&client_id=${configData.COGNITO_CLIENT_ID}&redirect_uri=${redirectUri}`;
    window.location.href = url;
  };

  public redirectSignup = () => {
    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/signup?response_type=code&client_id=${configData.COGNITO_CLIENT_ID}&redirect_uri=${redirectUri}`;
    window.location.href = url;
  };
}
