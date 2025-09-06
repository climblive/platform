<script lang="ts">
  import { ErrorBoundary, importNativeStyles } from "@climblive/lib/components";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import Start from "./failsafe/Start.svelte";

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });
</script>

<ErrorBoundary>
  {#await importNativeStyles() then styles}
    <svelte:component this={styles.default} />
  {/await}

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
