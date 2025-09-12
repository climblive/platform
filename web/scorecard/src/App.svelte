<script lang="ts">
  import EditProfile from "@/pages/EditProfile.svelte";
  import Loading from "@/pages/Loading.svelte";
  import Register from "@/pages/Register.svelte";
  import Scorecard from "@/pages/Scorecard.svelte";
  import Start from "@/pages/Start.svelte";
  import { type ScorecardSession } from "@/types";
  import { authenticateContender } from "@/utils/auth";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import { ErrorBoundary } from "@climblive/lib/components";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { onMount, setContext } from "svelte";
  import { Route, Router, navigate } from "svelte-routing";
  import { writable } from "svelte/store";
  import { ZodError } from "zod";

  let authenticating = $state(true);

  const session = writable<ScorecardSession>({
    contenderId: NaN,
    contestId: NaN,
    registrationCode: "",
    timestamp: new Date(0),
  });

  setContext("scorecardSession", session);

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  });

  const authenticate = async (code: string) => {
    try {
      const contender = await authenticateContender(code, queryClient, session);

      if (!contender.entered) {
        navigate(`/${code}/register`);
      }
    } catch (e) {
      if (e instanceof ZodError) {
        // eslint-disable-next-line no-console
        console.error(e);
      }

      navigate("/");
    } finally {
      authenticating = false;
    }
  };

  const parseCodeFromUrl = () => {
    const result = window.location.pathname.match(
      /^\/([0-9a-zA-Z]{8})(\/.*)?$/,
    );

    if (result) {
      authenticate(result[1]);
    } else {
      authenticating = false;
      navigate("/");
    }
  };

  parseCodeFromUrl();

  let compatibilityIgnored = $state(false);
  let code = $state<string>();

  const extractCodeFromPath = () => {
    const match = window.location.pathname.match(/\/([A-Z0-9]{8})/i);
    return match ? match[1] : null;
  };

  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    compatibilityIgnored = urlParams.get("compat") === "ignore";

    const extractedCode = extractCodeFromPath();

    if (extractedCode) {
      code = extractedCode;
    }
  });
</script>

<ErrorBoundary>
  <QueryClientProvider client={queryClient}>
    {#if compatibilityIgnored}
      <wa-callout variant="neutral" appearance="outlined filled">
        <wa-icon slot="icon" name="life-ring"></wa-icon>
        If you experience issues with the app you can
        <a href={code ? `/failsafe/${code}` : "/failsafe"}
          >switch to a basic version</a
        >.
      </wa-callout>
    {/if}
    {#if authenticating}
      <Loading />
    {:else}
      <Router>
        <Route path="/:code/register"><Register /></Route>
        <Route path="/:code/edit"><EditProfile /></Route>
        <Route path="/:code"><Scorecard /></Route>
        <Route path="/"><Start /></Route>
      </Router>
    {/if}
    {#if import.meta.env.DEV}
      <SvelteQueryDevtools />
    {/if}
  </QueryClientProvider>
</ErrorBoundary>

<style>
  wa-callout {
    margin: var(--wa-space-s);
    margin-block-end: 0;
  }
</style>
