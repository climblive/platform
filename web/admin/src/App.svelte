<script lang="ts">
  import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { Route, Router } from "svelte-routing";
  import Start from "@/pages/Start.svelte";

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
    <Route path="/"><Start /></Route>
  </Router>
  {#if import.meta.env.DEV && false}
    <SvelteQueryDevtools />
  {/if}
</QueryClientProvider>
