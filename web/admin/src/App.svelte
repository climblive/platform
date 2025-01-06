<script lang="ts">
  import Start from "@/pages/Start.svelte";
  import configData from "@climblive/lib/config.json";
  import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { Route, Router } from "svelte-routing";
  import Auth from "./pages/Auth.svelte";

  setBasePath("/shoelace");

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });

  const login = () => {
    const redirectUri = encodeURIComponent(
      window.location.origin + "/admin/auth",
    );
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/login?response_type=code&client_id=${configData.COGNITO_CLIENT_ID}&redirect_uri=${redirectUri}`;
    window.location.href = url;
  };
</script>

<QueryClientProvider client={queryClient}>
  <button onclick={login}>Login</button>
  <Router>
    <Route path="/admin"><Start /></Route>
    <Route path="/admin/auth"><Auth /></Route>
  </Router>
  {#if import.meta.env.DEV && false}
    <SvelteQueryDevtools />
  {/if}
</QueryClientProvider>
