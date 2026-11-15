<script lang="ts">
  import CompetitionDashboard from "@/components/CompetitionDashboard.svelte";
  import Loader from "@/components/Loader.svelte";
  import RulesEditor from "@/components/RulesEditor.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    getContendersByCompetitionQuery,
    getCompetitionQuery,
    getSelfQuery,
  } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";
  import ArchiveCompetition from "./ArchiveCompetition.svelte";
  import CompClassList from "./CompClassList.svelte";
  import DuplicateCompetition from "./DuplicateCompetition.svelte";
  import PopulateCompetition from "./PopulateCompetition.svelte";
  import ProblemList from "./ProblemList.svelte";
  import RaffleList from "./RaffleList.svelte";
  import RestoreCompetition from "./RestoreCompetition.svelte";
  import ScoreEngine from "./ScoreEngine.svelte";
  import TicketList from "./TicketList.svelte";
  import TransferCompetition from "./TransferCompetition.svelte";

  interface Props {
    competitionId: number;
  }

  let { competitionId }: Props = $props();

  let problemsHeading: HTMLHeadingElement | undefined = $state();
  let compClassesHeading: HTMLHeadingElement | undefined = $state();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const selfQuery = $derived(getSelfQuery());

  const competition = $derived(competitionQuery.data);
  const organizerId = $derived(competition?.ownership.organizerId);
  const self = $derived(selfQuery.data);

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

  const contendersQuery = $derived(
    getContendersByCompetitionQuery(competitionId),
  );
  const contenders = $derived(contendersQuery.data);

  const handleDownloadSimulatorConfig = () => {
    const registrationCodes = contenders?.map((c) => c.registrationCode) ?? [];

    const config = {
      apiUrl: new URL("/api", window.location.origin).toString(),
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
    a.download = `competition_${competitionId}_simulator_config.json`;
    a.click();

    window.URL.revokeObjectURL(url);
  };
</script>

<main>
  {#if competition === undefined || organizerId === undefined}
    <Loader />
  {:else}
    <wa-breadcrumb>
      <wa-breadcrumb-item
        onclick={() =>
          navigate(
            `/admin/organizers/${competition.ownership.organizerId}/competitions`,
          )}><wa-icon name="home"></wa-icon></wa-breadcrumb-item
      >
      <wa-breadcrumb-item>{competition.name}</wa-breadcrumb-item>
    </wa-breadcrumb>

    {#if competition.archivedAt != null}
      <RestoreCompetition {competitionId} />
    {:else}
      <CompetitionDashboard {competitionId} />

      <div class="results">
        <wa-button
          size="s"
          appearance="outlined"
          variant="brand"
          onclick={() =>
            navigate(`/admin/competitions/${competition.id}/results`)}
        >
          View results
          <wa-icon slot="end" name="arrow-right"></wa-icon>
        </wa-button>
      </div>

      <h2>Rules</h2>
      <wa-divider></wa-divider>
      <RulesEditor {competition} />

      <h2 bind:this={compClassesHeading}>Categories</h2>
      <wa-divider></wa-divider>
      <CompClassList {competitionId} />

      <h2 bind:this={problemsHeading}>Problems</h2>
      <wa-divider></wa-divider>
      <ProblemList
        {competitionId}
        {organizerId}
        tableLimit={window.location.hash.substring(1) === "problems"
          ? undefined
          : 8}
      />

      <h2>Tickets</h2>
      <wa-divider></wa-divider>
      <TicketList {competitionId} />

      <h2>Raffles</h2>
      <wa-divider></wa-divider>
      <RaffleList {competitionId} />

      <h2>Advanced</h2>
      <wa-divider></wa-divider>
      <h3>Actions</h3>
      <div class="actions">
        <DuplicateCompetition
          {competitionId}
          competitionName={competition.name}
        >
          {#snippet children({ duplicateCompetition })}
            <wa-button onclick={duplicateCompetition} appearance="outlined">
              Duplicate
              <wa-icon name="copy" slot="start"></wa-icon>
            </wa-button>
          {/snippet}
        </DuplicateCompetition>
        <TransferCompetition
          {competitionId}
          {organizerId}
          competitionName={competition.name}
        >
          {#snippet children({ transferCompetition, disabled })}
            <wa-button
              onclick={transferCompetition}
              appearance="outlined"
              {disabled}
            >
              Transfer
              <wa-icon name="arrow-right" slot="start"></wa-icon>
            </wa-button>
          {/snippet}
        </TransferCompetition>
        <ArchiveCompetition
          {competitionId}
          {organizerId}
          competitionName={competition.name}
        >
          {#snippet children({ archiveCompetition })}
            <wa-button
              onclick={archiveCompetition}
              appearance="outlined"
              variant="danger"
            >
              Archive
              <wa-icon name="box-archive" slot="start"></wa-icon>
            </wa-button>
          {/snippet}
        </ArchiveCompetition>
      </div>
      {#if location.hostname !== "climblive.app" || self?.admin}
        <h3>Developer tools</h3>
        <div class="developer-tools">
          <wa-button
            appearance="outlined"
            disabled={!contenders || contenders.length === 0}
            onclick={handleDownloadSimulatorConfig}
            >Download simulator config
            <wa-icon name="download" slot="start"></wa-icon>
          </wa-button>

          <PopulateCompetition {competitionId} />
        </div>
      {/if}
      <h3>Score Engines</h3>
      <p>
        An active score engine collects all results during a competition and
        computes scores and rankings for all participants.
      </p>

      <ScoreEngine {competitionId} />
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

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
  }

  .developer-tools {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
  }

  .results {
    margin-block-start: var(--wa-space-m);
  }
</style>
