<script lang="ts">
  import { ErrorBoundary } from "@climblive/lib/components";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import Start from "./failsafe/Start.svelte";

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });

  const handleError = (event: ErrorEvent) => {
    alert(event.message);
  };
</script>

<svelte:window onerror={handleError} />

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
