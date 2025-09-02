<script lang="ts">
  import { ErrorBoundary } from "@climblive/lib/components";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { onMount } from "svelte";
  import Start from "./failsafe/Start.svelte";

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });

  onMount(async () => {
    if (false) {
      const NativeStyles = (await import("./NativeStyles.svelte")).default;

      new NativeStyles({
        target: document.body,
      });
    }
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
