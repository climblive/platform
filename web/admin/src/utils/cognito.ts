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
    Authorization: `Basic ${btoa(configData.COGNITO_CLIENT_ID + ":" + configData.COGNITO_CLIENT_SECRET)}`
  },
});

export const exchangeCode = async (code: string) => {
  const params = new URLSearchParams();
  params.append("grant_type", "authorization_code");
  params.append("client_id", configData.COGNITO_CLIENT_ID);
  params.append("code", code);
  params.append(
    "redirect_uri",
    window.location.protocol + "//" + window.location.host + "/admin",
  );

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
