<script lang="ts">
  import "@shoelace-style/shoelace/dist/themes/light.css";
  import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { setContext } from "svelte";
  import { Route, Router, navigate } from "svelte-routing";
  import { writable } from "svelte/store";
  import "./main.css";
  import EditProfile from "./pages/EditProfile.svelte";
  import Loading from "./pages/Loading.svelte";
  import Register from "./pages/Register.svelte";
  import Scorecard from "./pages/Scorecard.svelte";
  import Start from "./pages/Start.svelte";
  import { type ScorecardSession } from "./types";
  import { authenticateContender } from "./utils/auth";

  let authenticating = true;

  setBasePath("/shoelace");

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
      navigate("/");
    } finally {
      authenticating = false;
    }
  };

  const parseCodeFromUrl = () => {
    const result = window.location.pathname.match(
      /^\/([0-9a-zA-Z]{8})(\/.*)?$/
    );

    if (result) {
      authenticate(result[1]);
    } else {
      authenticating = false;
      navigate("/");
    }
  };

  parseCodeFromUrl();
</script>

<QueryClientProvider client={queryClient}>
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
  {#if import.meta.env.DEV && false}
    <SvelteQueryDevtools />
  {/if}
</QueryClientProvider>
