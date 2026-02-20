<script lang="ts">
  import ContestDashboard from "@/components/ContestDashboard.svelte";
  import Loader from "@/components/Loader.svelte";
  import RulesEditor from "@/components/RulesEditor.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { getContestQuery } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";
  import ArchiveContest from "./ArchiveContest.svelte";
  import CompClassList from "./CompClassList.svelte";
  import DuplicateContest from "./DuplicateContest.svelte";
  import ProblemList from "./ProblemList.svelte";
  import RaffleList from "./RaffleList.svelte";
  import RestoreContest from "./RestoreContest.svelte";
  import ScoreEngine from "./ScoreEngine.svelte";
  import TicketList from "./TicketList.svelte";
  import TransferContest from "./TransferContest.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let problemsHeading: HTMLHeadingElement | undefined = $state();
  let compClassesHeading: HTMLHeadingElement | undefined = $state();

  const contestQuery = $derived(getContestQuery(contestId));

  const contest = $derived(contestQuery.data);

  $effect(() => {
    const hash = window.location.hash.substring(1);

    let scrollElement: HTMLElement | undefined;

    switch (hash) {
      case "problems":
        scrollElement = problemsHeading;
        break;
      case "comp-classes":
        scrollElement = compClassesHeading;
        break;
    }

    if (scrollElement === undefined) {
      return;
    }

    const cb = () =>
      scrollElement.scrollIntoView({
        behavior: "instant",
        block: "start",
        inline: "nearest",
      });

    if (window.requestIdleCallback !== undefined) {
      requestIdleCallback(cb);
    } else {
      setTimeout(cb, 250);
    }
  });
</script>

<main>
  {#if contest === undefined}
    <Loader />
  {:else}
    <wa-breadcrumb>
      <wa-breadcrumb-item
        onclick={() =>
          navigate(
            `/admin/organizers/${contest.ownership.organizerId}/contests`,
          )}><wa-icon name="home"></wa-icon></wa-breadcrumb-item
      >
      <wa-breadcrumb-item>{contest.name}</wa-breadcrumb-item>
    </wa-breadcrumb>

    {#if contest.archived === true}
      <RestoreContest {contestId} />
    {/if}

    {#if contest.archived === false}
      <ContestDashboard contestId={contest.id} />
      <div class="results-section">
        <wa-button
          size="small"
          class="view-results-button"
          appearance="outlined"
          variant="brand"
          onclick={() => navigate(`/admin/contests/${contest.id}/results`)}
        >
          View results
          <wa-icon slot="end" name="arrow-right"></wa-icon>
        </wa-button>
      </div>

      <h2>Rules</h2>
      <wa-divider></wa-divider>
      <RulesEditor {contest} />

      <h2 bind:this={compClassesHeading}>Classes</h2>
      <wa-divider></wa-divider>
      <CompClassList {contestId} />

      <h2>Tickets</h2>
      <wa-divider></wa-divider>
      <TicketList {contestId} />

      <h2 bind:this={problemsHeading}>Problems</h2>
      <wa-divider></wa-divider>
      <ProblemList
        {contestId}
        organizerId={contest.ownership.organizerId}
        tableLimit={window.location.hash.substring(1) === "problems"
          ? undefined
          : 8}
      />

      <h2>Raffles</h2>
      <wa-divider></wa-divider>
      <RaffleList {contestId} />

      <h2>Advanced</h2>
      <wa-divider></wa-divider>
      <h3>Actions</h3>
      <div class="actions">
        <DuplicateContest {contestId} />
        <TransferContest
          {contestId}
          organizerId={contest.ownership.organizerId}
        />
        <ArchiveContest
          {contestId}
          organizerId={contest.ownership.organizerId}
        />
      </div>
      <h3>Score Engines</h3>
      <p>
        An active score engine collects all results during a contest and
        computes scores and rankings for all participants.
      </p>

      <ScoreEngine {contestId} />
    {/if}
  {/if}
</main>

<style>
  wa-breadcrumb {
    margin-block-end: var(--wa-space-m);
    display: block;
  }

  h2 {
    margin-top: var(--wa-space-2xl);
  }

  wa-divider {
    --color: var(--wa-color-brand-fill-normal);
  }

  .results-section {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-m);
    flex-wrap: wrap;
    margin-block-start: var(--wa-space-l);
  }

  .view-results-button {
    margin-inline-start: auto;
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
  }
</style>
