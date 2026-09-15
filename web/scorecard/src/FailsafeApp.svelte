<script lang="ts">
  import { ErrorBoundary } from "@climblive/lib/components";
  import { HOUR } from "@climblive/lib/queries";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import Start from "./failsafe/Start.svelte";

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 12 * HOUR,
        gcTime: 12 * HOUR,
        refetchOnWindowFocus: "always",
      },
    },
  });
</script>

<ErrorBoundary>
  <QueryClientProvider client={queryClient}>
    <main>
      <Start />
    </main>
  </QueryClientProvider>
</ErrorBoundary>

<style>
  main {
    padding: var(--wa-space-m);
  }
</style>
