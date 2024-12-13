<script lang="ts">
  import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { Route, Router } from "svelte-routing";
  import Scoreboard from "./pages/Scoreboard.svelte";

  setBasePath("/shoelace");

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });
</script>

<QueryClientProvider client={queryClient}>
  <Router>
    <Route path="/scoreboard/:contestId"
      >{#snippet children({ params })}
        <Scoreboard contestId={Number(params.contestId)} />
      {/snippet}
    </Route>
  </Router>
  {#if import.meta.env.DEV && false}
    <SvelteQueryDevtools />
  {/if}
</QueryClientProvider>
