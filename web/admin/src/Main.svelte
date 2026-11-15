<script lang="ts">
  import { Route, Router } from "svelte-routing";
  import Header from "./Header.svelte";
  import CompetitorView from "./pages/CompetitorView.svelte";
  import Competition from "./pages/Competition.svelte";
  import CreateCategory from "./pages/CreateCategory.svelte";
  import CreateCompetition from "./pages/CreateCompetition.svelte";
  import CreateProblem from "./pages/CreateProblem.svelte";
  import EditCategory from "./pages/EditCategory.svelte";
  import EditCompetition from "./pages/EditCompetition.svelte";
  import EditProblem from "./pages/EditProblem.svelte";
  import Health from "./pages/Health.svelte";
  import Help from "./pages/Help.svelte";
  import InviteView from "./pages/InviteView.svelte";
  import ManageOrganizer from "./pages/ManageOrganizer.svelte";
  import OrganizerView from "./pages/OrganizerView.svelte";
  import PrintableTicketList from "./pages/PrintableTicketList.svelte";
  import Profile from "./pages/Profile.svelte";
  import RaffleView from "./pages/RaffleView.svelte";
  import ResultsPage from "./pages/ResultsPage.svelte";
  import Root from "./pages/Root.svelte";
  import TicketsListPage from "./pages/TicketsListPage.svelte";
</script>

<Header />

<main>
  <Router basepath="/admin">
    <Route path="/">
      <Root />
    </Route>
    <Route path="/organizers/:organizerId">
      {#snippet children({ params }: { params: { organizerId: number } })}
        <ManageOrganizer organizerId={Number(params.organizerId)} />
      {/snippet}
    </Route>
    <Route path="/invites/:inviteId">
      {#snippet children({ params }: { params: { inviteId: string } })}
        <InviteView inviteId={params.inviteId} />
      {/snippet}
    </Route>
    <Route path="/help">
      <Help />
    </Route>
    <Route path="/profile">
      <Profile />
    </Route>
    <Route path="/health">
      <Health />
    </Route>
    <Route path="/organizers/:organizerId/competitions">
      {#snippet children({ params }: { params: { organizerId: number } })}
        <OrganizerView organizerId={Number(params.organizerId)} />
      {/snippet}
    </Route>
    <Route path="/organizers/:organizerId/competitions/new">
      {#snippet children({ params }: { params: { organizerId: number } })}
        <CreateCompetition organizerId={Number(params.organizerId)} />
      {/snippet}
    </Route>
    <Route path="/competitions/:competitionId">
      {#snippet children({ params }: { params: { competitionId: number } })}
        {#key params.competitionId}
          <Competition competitionId={Number(params.competitionId)} />
        {/key}
      {/snippet}
    </Route>
    <Route path="/competitions/:competitionId/edit">
      {#snippet children({ params }: { params: { competitionId: number } })}
        <EditCompetition competitionId={Number(params.competitionId)} />
      {/snippet}
    </Route>
    <Route path="/competitions/:competitionId/new-category">
      {#snippet children({ params }: { params: { competitionId: number } })}
        <CreateCategory competitionId={Number(params.competitionId)} />
      {/snippet}
    </Route>
    <Route path="/competitions/:competitionId/new-problem">
      {#snippet children({ params }: { params: { competitionId: number } })}
        <CreateProblem competitionId={Number(params.competitionId)} />
      {/snippet}
    </Route>
    <Route path="/problems/:problemId/edit">
      {#snippet children({ params }: { params: { problemId: number } })}
        <EditProblem problemId={Number(params.problemId)} />
      {/snippet}
    </Route>
    <Route path="/categories/:categoryId/edit">
      {#snippet children({ params }: { params: { categoryId: number } })}
        <EditCategory categoryId={Number(params.categoryId)} />
      {/snippet}
    </Route>
    <Route path="/raffles/:raffleId">
      {#snippet children({ params }: { params: { raffleId: number } })}
        <RaffleView raffleId={Number(params.raffleId)} />
      {/snippet}
    </Route>
    <Route path="/competitions/:competitionId/results">
      {#snippet children({ params }: { params: { competitionId: number } })}
        <ResultsPage competitionId={Number(params.competitionId)} />
      {/snippet}
    </Route>
    <Route path="/competitions/:competitionId/tickets">
      {#snippet children({ params }: { params: { competitionId: number } })}
        <TicketsListPage competitionId={Number(params.competitionId)} />
      {/snippet}
    </Route>
    <Route path="/competitions/:competitionId/tickets/print">
      {#snippet children({ params }: { params: { competitionId: number } })}
        <PrintableTicketList competitionId={Number(params.competitionId)} />
      {/snippet}
    </Route>
    <Route path="/competitors/:competitorId">
      {#snippet children({ params }: { params: { competitorId: number } })}
        <CompetitorView competitorId={Number(params.competitorId)} />
      {/snippet}
    </Route>
  </Router>
</main>

<style>
  main {
    padding: var(--wa-space-m);
    margin: 0 auto;
    max-width: 1024px;
  }

  @media print {
    main {
      padding: 0;
    }
  }
</style>
