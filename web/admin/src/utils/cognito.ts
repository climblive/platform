import configData from "@climblive/lib/config.json";
import axios from "axios";

interface OAuthTokenResponse {
  id_token: string;
  access_token: string;
  refresh_token: string;
  expires_id: number;
  token_type: string;
}

const instance = axios.create({
  baseURL: "https://clmb.auth.eu-west-1.amazoncognito.com",
  timeout: 10_000,
  headers: {
    "Content-Type": "application/x-www-form-urlencoded",
  },
});

export const exchangeCode = async (code: string, state: string) => {
  const codeVerifier = sessionStorage.getItem("oauth_code_verifier");
  const expectedState = sessionStorage.getItem("oauth_state");
  sessionStorage.removeItem("oauth_code_verifier");
  sessionStorage.removeItem("oauth_state");

  if (
    codeVerifier === null ||
    expectedState === null ||
    state !== expectedState
  ) {
    throw new Error("Invalid OAuth response");
  }

  const params = new URLSearchParams();
  params.append("grant_type", "authorization_code");
  params.append("client_id", configData.COGNITO_CLIENT_ID);
  params.append("code", code);
  params.append(
    "redirect_uri",
    window.location.protocol + "//" + window.location.host + "/admin",
  );

  params.append("oauth_code_verifier", codeVerifier);

  const response = await instance.post("/oauth2/token", params);

  return response.data as OAuthTokenResponse;
};

export const refreshSession = async (refreshToken: string) => {
  const params = new URLSearchParams();
  params.append("grant_type", "refresh_token");
  params.append("client_id", configData.COGNITO_CLIENT_ID);
  params.append("refresh_token", refreshToken);

  const response = await instance.post("/oauth2/token", params);

  return response.data as Omit<OAuthTokenResponse, "refresh_token">;
};
