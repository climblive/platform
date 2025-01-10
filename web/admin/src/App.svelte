<script lang="ts">
  import { ApiClient, OrganizerCredentialsProvider } from "@climblive/lib";
  import configData from "@climblive/lib/config.json";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/spinner/spinner.js";
  import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { navigate, Route, Router } from "svelte-routing";
  import Contest from "./pages/Contest.svelte";
  import { exchangeCode, refreshSession } from "./utils/cognito";

  setBasePath("/shoelace");

  let authenticated = $state(false);

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });

  const authenticate = async () => {
    const query = new URLSearchParams(location.search);
    const code = query.get("code");

    if (code != null) {
      const { access_token, refresh_token } = await exchangeCode(code);

      ApiClient.getInstance().setCredentialsProvider(
        new OrganizerCredentialsProvider(access_token),
      );

      localStorage.setItem("refresh_token", refresh_token);

      authenticated = true;

      navigate("/");

      return;
    }

    try {
      const refreshToken = localStorage.getItem("refresh_token");

      if (!refreshToken) {
        return;
      }

      if (refreshToken) {
        const { access_token } = await refreshSession(refreshToken);

        ApiClient.getInstance().setCredentialsProvider(
          new OrganizerCredentialsProvider(access_token),
        );

        authenticated = true;
      }
    } catch {
      localStorage.removeItem("refresh_token");
    }
  };

  const login = () => {
    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/login?response_type=code&client_id=${configData.COGNITO_CLIENT_ID}&redirect_uri=${redirectUri}`;
    window.location.href = url;
  };
</script>

{#await authenticate()}
  <sl-spinner></sl-spinner>
{:then}
  <QueryClientProvider client={queryClient}>
    {#if !authenticated}
      <sl-button variant="primary" onclick={login}>Login</sl-button>
    {/if}
    <Router basepath="/admin">
      <Route path="/contests/:contestId">
        {#snippet children({ params })}
          <Contest contestId={Number(params.contestId)} />
        {/snippet}
      </Route>
    </Router>
    {#if import.meta.env.DEV && false}
      <SvelteQueryDevtools />
    {/if}
  </QueryClientProvider>
{/await}
