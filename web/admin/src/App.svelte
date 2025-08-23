<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/spinner/spinner.js";
  import { ErrorBoundary } from "@climblive/lib/components";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { onDestroy, onMount, setContext } from "svelte";
  import { navigate } from "svelte-routing";
  import { writable } from "svelte/store";
  import { Authenticator } from "./authenticator.svelte";
  import Main from "./Main.svelte";

  const selectedOrganizer = writable<number | undefined>();
  const authenticator = new Authenticator();

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

  const handleVisibilityChange = () => {
    if (document.visibilityState === "visible") {
      authenticator.startKeepAlive();
    } else {
      authenticator.stopKeepAlive();
    }
  };

  onMount(() => {
    const regex = /^\/admin\/organizers\/(\d+)$/;

    const match = window.location.pathname.match(regex);

    if (match) {
      const organizerId = Number(match[1]);
      $selectedOrganizer = organizerId;
    }
  });

  onMount(authenticator.startKeepAlive);

  onDestroy(authenticator.stopKeepAlive);
</script>

<svelte:window
  onstorage={handleStorageEvent}
  onvisibilitychange={handleVisibilityChange}
/>

<ErrorBoundary>
  {#await authenticator.authenticate()}
    <main>
      <wa-spinner></wa-spinner>
    </main>
  {:then}
    <QueryClientProvider client={queryClient}>
      {#if !authenticator.isAuthenticated()}
        <main>
          <section>
            <h1>Hi!</h1>
            <p>
              Welcome to the <em>brand new</em> admin console for ClimbLive.
            </p>
            <wa-button variant="brand" onclick={authenticator.redirectLogin}
              >Login</wa-button
            >
            <wa-button
              variant="brand"
              appearance="plain"
              onclick={authenticator.redirectSignup}>Sign up</wa-button
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
    height: 100vh;
    padding: var(--wa-space-l);
    padding-top: 20vh;
  }

  wa-spinner {
    font-size: 5rem;
  }
</style>
