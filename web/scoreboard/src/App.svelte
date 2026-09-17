<script lang="ts">
  import "@awesome.me/webawesome/dist/components/toast/toast.js";
  import { ErrorBoundary } from "@climblive/lib/components";
  import { HOUR, MINUTE } from "@climblive/lib/queries";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { Route, Router } from "svelte-routing";
  import Scoreboard from "./pages/Scoreboard.svelte";

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 5 * MINUTE,
        gcTime: 12 * HOUR,
        refetchOnWindowFocus: true,
      },
    },
  });
</script>

<ErrorBoundary>
  <wa-toast placement="bottom-end"></wa-toast>

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
