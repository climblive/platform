<script lang="ts">
  import { ApiClient, OrganizerCredentialsProvider } from "@climblive/lib";
  import { ErrorBoundary } from "@climblive/lib/components";
  import configData from "@climblive/lib/config.json";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/spinner/spinner.js";
  import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { navigate, Route, Router } from "svelte-routing";
  import Contest from "./pages/Contest.svelte";
  import ContestList from "./pages/ContestList.svelte";
  import CreateCompClass from "./pages/CreateCompClass.svelte";
  import CreateContest from "./pages/CreateContest.svelte";
  import CreateProblem from "./pages/CreateProblem.svelte";
  import EditCompClass from "./pages/EditCompClass.svelte";
  import EditProblem from "./pages/EditProblem.svelte";
  import { exchangeCode, refreshSession } from "./utils/cognito";

  setBasePath("/shoelace");

  let authenticated = $state(false);

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: true,
      },
    },
  });

  const authenticate = async () => {
    const query = new URLSearchParams(location.search);
    const code = query.get("code");

    if (code != null) {
      const { access_token, refresh_token } = await exchangeCode(code);

      ApiClient.getInstance().setCredentialsProvider(
        new OrganizerCredentialsProvider(access_token),
      );

      localStorage.setItem("refresh_token", refresh_token);

      authenticated = true;

      navigate("./", { replace: true });

      return;
    }

    try {
      const refreshToken = localStorage.getItem("refresh_token");

      if (!refreshToken) {
        return;
      }

      if (refreshToken) {
        const { access_token } = await refreshSession(refreshToken);

        ApiClient.getInstance().setCredentialsProvider(
          new OrganizerCredentialsProvider(access_token),
        );

        authenticated = true;
      }
    } catch {
      localStorage.removeItem("refresh_token");
    }
  };

  const login = () => {
    const redirectUri = encodeURIComponent(window.location.origin + "/admin");
    const url = `https://clmb.auth.eu-west-1.amazoncognito.com/login?response_type=code&client_id=${configData.COGNITO_CLIENT_ID}&redirect_uri=${redirectUri}`;
    window.location.href = url;
  };
</script>

<ErrorBoundary>
  {#await authenticate()}
    <sl-spinner></sl-spinner>
  {:then}
    <QueryClientProvider client={queryClient}>
      {#if !authenticated}
        <sl-button variant="primary" onclick={login}>Login</sl-button>
      {/if}
      <main>
        <Router basepath="/admin">
          <Route path="/organizers/:organizerId">
            {#snippet children({ params }: { params: { organizerId: number } })}
              <ContestList organizerId={Number(params.organizerId)} />
            {/snippet}
          </Route>
          <Route path="/organizers/:organizerId/contests/new">
            {#snippet children({ params }: { params: { organizerId: number } })}
              <CreateContest organizerId={Number(params.organizerId)} />
            {/snippet}
          </Route>
          <Route path="/contests/:contestId">
            {#snippet children({ params }: { params: { contestId: number } })}
              <Contest contestId={Number(params.contestId)} />
            {/snippet}
          </Route>
          <Route path="/contests/:contestId/new-comp-class">
            {#snippet children({ params }: { params: { contestId: number } })}
              <CreateCompClass contestId={Number(params.contestId)} />
            {/snippet}
          </Route>
          <Route path="/contests/:contestId/new-problem">
            {#snippet children({ params }: { params: { contestId: number } })}
              <CreateProblem contestId={Number(params.contestId)} />
            {/snippet}
          </Route>
          <Route path="/problems/:problemId/edit">
            {#snippet children({ params }: { params: { problemId: number } })}
              <EditProblem problemId={Number(params.problemId)} />
            {/snippet}
          </Route>
          <Route path="/comp-classes/:compClassId/edit">
            {#snippet children({ params }: { params: { compClassId: number } })}
              <EditCompClass compClassId={Number(params.compClassId)} />
            {/snippet}
          </Route>
        </Router>
      </main>
      {#if import.meta.env.DEV}
        <SvelteQueryDevtools />
      {/if}
    </QueryClientProvider>
  {/await}
</ErrorBoundary>

<style>
  main {
    padding: var(--sl-spacing-medium);
  }
</style>
