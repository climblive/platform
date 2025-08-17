<script lang="ts">
  import "@awesome.me/webawesome/dist/components/spinner/spinner.js";
  import { Route, Router } from "svelte-routing";
  import Header from "./Header.svelte";
  import Contest from "./pages/Contest.svelte";
  import CreateCompClass from "./pages/CreateCompClass.svelte";
  import CreateContest from "./pages/CreateContest.svelte";
  import CreateProblem from "./pages/CreateProblem.svelte";
  import EditCompClass from "./pages/EditCompClass.svelte";
  import EditProblem from "./pages/EditProblem.svelte";
  import OrganizerView from "./pages/OrganizerView.svelte";
  import PrintableTicketList from "./pages/PrintableTicketList.svelte";
  import RaffleView from "./pages/RaffleView.svelte";
  import Root from "./pages/Root.svelte";
</script>

<Header />

<main>
  <Router basepath="/admin">
    <Route path="/">
      <Root />
    </Route>
    <Route path="/organizers/:organizerId">
      {#snippet children({ params }: { params: { organizerId: number } })}
        <OrganizerView organizerId={Number(params.organizerId)} />
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
    <Route path="/raffles/:raffleId">
      {#snippet children({ params }: { params: { raffleId: number } })}
        <RaffleView raffleId={Number(params.raffleId)} />
      {/snippet}
    </Route>
    <Route path="/contests/:contestId/tickets">
      {#snippet children({ params }: { params: { contestId: number } })}
        <PrintableTicketList contestId={Number(params.contestId)} />
      {/snippet}
    </Route>
  </Router>
</main>

<style>
  main {
    padding: var(--wa-space-m);
  }

  @media print {
    main {
      padding: 0;
    }
  }
</style>
