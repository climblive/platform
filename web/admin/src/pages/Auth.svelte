<script lang="ts">
  import { ApiClient, OrganizerCredentialsProvider } from "@climblive/lib";
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
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
  });

  const query = new URLSearchParams(location.search);
  const code = query.get("code");

  if (code != null) {
    const params = new URLSearchParams();
    params.append("grant_type", "authorization_code");
    params.append("client_id", configData.COGNITO_CLIENT_ID);
    params.append("code", code);
    params.append(
      "redirect_uri",
      window.location.protocol + "//" + window.location.host + "/admin/auth",
    );

    instance
      .post("/oauth2/token", params, {
        headers: {
          Authorization: `Basic ${btoa(configData.COGNITO_CLIENT_ID + ":" + configData.COGNITO_CLIENT_SECRET)}`,
        },
      })
      .then(function (response) {
        const { id_token, access_token, refresh_token }: OAuthTokenResponse =
          response.data;

        ApiClient.getInstance().setCredentialsProvider(
          new OrganizerCredentialsProvider(access_token),
        );

        ApiClient.getInstance().getTicks(1);
      })
      .catch(function (error) {});
  }
</script>
