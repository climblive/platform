<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/spinner/spinner.js";
  import { ApiClient, OrganizerCredentialsProvider } from "@climblive/lib";
  import { ErrorBoundary } from "@climblive/lib/components";
  import configData from "@climblive/lib/config.json";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { onMount, setContext } from "svelte";
  import { navigate } from "svelte-routing";
  import { writable } from "svelte/store";
  import Main from "./Main.svelte";
  import { exchangeCode, refreshSession } from "./utils/cognito";

  let authenticated = $state(false);

  const selectedOrganizer = writable<number | undefined>();

  const organizerId = localStorage.getItem("organizerId");
  if (organizerId !== null) {
    $selectedOrganizer = Number(organizerId);
  }

  $effect(() => {
    if ($selectedOrganizer !== undefined) {
      localStorage.setItem("organizerId", $selectedOrganizer.toString());
    }
  });

  const handleStorageEvent = (e: StorageEvent) => {
    if (e.key !== "organizerId") {
      return;
    }

    $selectedOrganizer = Number(e.newValue);
    navigate(`/admin/organizers/${e.newValue}`);
  };

  setContext("selectedOrganizer", selectedOrganizer);

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: true,
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

      navigate("./", { replace: true });

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

  const signup = () => {
    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/signup?response_type=code&client_id=${configData.COGNITO_CLIENT_ID}&redirect_uri=${redirectUri}`;
    window.location.href = url;
  };

  onMount(() => {
    const regex = /^\/admin\/organizers\/(\d+)$/;

    const match = window.location.pathname.match(regex);

    if (match) {
      const organizerId = Number(match[1]);
      $selectedOrganizer = organizerId;
    }
  });
</script>

<svelte:window onstorage={handleStorageEvent} />

<ErrorBoundary>
  {#await authenticate()}
    <main>
      <wa-spinner></wa-spinner>
    </main>
  {:then}
    <QueryClientProvider client={queryClient}>
      {#if !authenticated}
        <main>
          <section>
            <h1>Hi!</h1>
            <p>
              Welcome to the <em>brand new</em> admin console for ClimbLive.
            </p>
            <wa-button variant="brand" onclick={login}>Login</wa-button>
            <wa-button variant="brand" appearance="plain" onclick={signup}
              >Sign up</wa-button
            >
          </section>
        </main>
      {:else}
        <Main />
      {/if}
      {#if import.meta.env.DEV}
        <SvelteQueryDevtools />
      {/if}
    </QueryClientProvider>
  {/await}
</ErrorBoundary>

<style>
  main {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
  }

  wa-spinner {
    font-size: 5rem;
  }
</style>
