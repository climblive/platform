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
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import {
    getContendersByContestQuery,
    getContestQuery,
  } from "@climblive/lib/queries";
  import { getApiUrl } from "@climblive/lib/utils";
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
  const sectionHelpId = $props.id();

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

  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const contenders = $derived(contendersQuery.data);

  const handleDownloadSimulatorConfig = () => {
    const registrationCodes = contenders?.map((c) => c.registrationCode) ?? [];

    const config = {
      apiUrl: new URL(getApiUrl(), window.location.origin).toString(),
      registrationCodes,
      iterations: 10,
      maxSleep: 10_000,
    };

    const blob = new Blob([JSON.stringify(config, null, 2)], {
      type: "application/json",
    });

    const url = window.URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.href = url;
    a.download = `contest_${contestId}_simulator_config.json`;
    a.click();

    window.URL.revokeObjectURL(url);
  };
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

    {#if contest.archivedAt != null}
      <RestoreContest {contestId} />
    {:else}
      <ContestDashboard contestId={contest.id} />

      <div class="results">
        <wa-button
          size="s"
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

      <h2 class="section-heading" bind:this={compClassesHeading}>
        Classes
        <wa-button
          id={`${sectionHelpId}-comp-classes`}
          size="s"
          appearance="plain"
          variant="neutral"
          aria-label="About classes"
        >
          <wa-icon name="lightbulb" variant="regular"></wa-icon>
        </wa-button>
      </h2>
      <wa-popover for={`${sectionHelpId}-comp-classes`}>
        Classes represent the categories in which the contenders compete,
        typically divided into Males and Females. The contest duration is
        defined by the start and end times of your classes.
        <wa-button data-popover="close" variant="primary" size="s"
          >Got it!</wa-button
        >
      </wa-popover>
      <wa-divider></wa-divider>
      <CompClassList {contestId} />

      <h2 class="section-heading" bind:this={problemsHeading}>
        Problems
        <wa-button
          id={`${sectionHelpId}-problems`}
          size="s"
          appearance="plain"
          variant="neutral"
          aria-label="About problems"
        >
          <wa-icon name="lightbulb" variant="regular"></wa-icon>
        </wa-button>
      </h2>
      <wa-popover for={`${sectionHelpId}-problems`}>
        Problems refer to the boulder problems that the contenders will attempt
        during the contest, each of which can have its own point value.
        <wa-button data-popover="close" variant="primary" size="s"
          >Got it!</wa-button
        >
      </wa-popover>
      <wa-divider></wa-divider>
      <ProblemList
        {contestId}
        organizerId={contest.ownership.organizerId}
        tableLimit={window.location.hash.substring(1) === "problems"
          ? undefined
          : 8}
      />

      <h2>Tickets</h2>
      <wa-divider></wa-divider>
      <TicketList {contestId} />

      <h2 class="section-heading">
        Raffles
        <wa-button
          id={`${sectionHelpId}-raffles`}
          size="s"
          appearance="plain"
          variant="neutral"
          aria-label="About raffles"
        >
          <wa-icon name="lightbulb" variant="regular"></wa-icon>
        </wa-button>
      </h2>
      <wa-popover for={`${sectionHelpId}-raffles`}>
        Raffles are used to randomly select prize winners, typically after the
        contest has ended.
        <wa-button data-popover="close" variant="primary" size="s"
          >Got it!</wa-button
        >
      </wa-popover>
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
      {#if location.hostname !== "climblive.app"}
        <h3>Developer tools</h3>
        <wa-button
          appearance="outlined"
          disabled={!contenders || contenders.length === 0}
          onclick={handleDownloadSimulatorConfig}
          >Download simulator config
          <wa-icon name="download" slot="start"></wa-icon>
        </wa-button>
      {/if}
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

  .section-heading {
    display: flex;
    align-items: center;
    gap: var(--wa-space-2xs);
  }

  wa-divider {
    --color: var(--wa-color-brand-fill-normal);
  }

  wa-popover {
    --max-width: 300px;
  }

  wa-popover wa-button {
    margin-block-start: var(--wa-space-m);
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
  }

  .results {
    margin-block-start: var(--wa-space-m);
  }
</style>
