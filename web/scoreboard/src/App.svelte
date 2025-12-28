<script lang="ts">
  import { ErrorBoundary, SplashScreen } from "@climblive/lib/components";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { Route, Router } from "svelte-routing";
  import Scoreboard from "./pages/Scoreboard.svelte";

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });

  let showSplash = $state(true);
</script>

{#if showSplash}
  <SplashScreen onComplete={() => (showSplash = false)} />
{:else}
  <ErrorBoundary>
    <QueryClientProvider client={queryClient}>
      <Router>
        <Route path="/scoreboard/:contestId"
          >{#snippet children({ params }: { params: { contestId: number } })}
            <Scoreboard contestId={Number(params.contestId)} />
          {/snippet}
        </Route>
      </Router>
      {#if import.meta.env.DEV}
        <SvelteQueryDevtools />
      {/if}
    </QueryClientProvider>
  </ErrorBoundary>
{/if}
