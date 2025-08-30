<script lang="ts">
  import { ErrorBoundary } from "@climblive/lib/components";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import NativeStyles from "./NativeStyles.svelte";
  import Start from "./fallback/Start.svelte";

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });
</script>

<NativeStyles />

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
