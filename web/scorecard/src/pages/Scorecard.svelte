<script lang="ts">
  import Header from "@/components/Header.svelte";
  import ProblemView from "@/components/ProblemView.svelte";
  import type { ScorecardSession } from "@/types";
  import type { ContestState } from "@/types/state";
  import { ResultList, ScoreboardProvider } from "@climblive/lib/components";
  import configData from "@climblive/lib/config.json";
  import type { ContenderScoreUpdatedEvent } from "@climblive/lib/models";
  import {
    getCompClassesQuery,
    getContenderQuery,
    getContestQuery,
    getProblemsQuery,
    getTicksQuery,
  } from "@climblive/lib/queries";
  import type { SlTabGroup, SlTabShowEvent } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/tab-group/tab-group.js";
  import "@shoelace-style/shoelace/dist/components/tab-panel/tab-panel.js";
  import "@shoelace-style/shoelace/dist/components/tab/tab.js";
  import { parseISO } from "date-fns";
  import { add } from "date-fns/add";
  import { differenceInMilliseconds } from "date-fns/differenceInMilliseconds";
  import { isBefore } from "date-fns/isBefore";
  import { getContext, onDestroy, onMount } from "svelte";
  import { type Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const contenderQuery = getContenderQuery($session.contenderId);
  const contestQuery = getContestQuery($session.contestId);
  const compClassesQuery = getCompClassesQuery($session.contestId);
  const problemsQuery = getProblemsQuery($session.contestId);
  const ticksQuery = getTicksQuery($session.contenderId);

  let resultsConnected = false;
  let tabGroup: SlTabGroup;
  let eventSource: EventSource | undefined;
  let score: number;
  let placement: number | undefined;
  let updateStateTimerId: number = 0;
  let state: ContestState = "NOT_STARTED";

  $: contender = $contenderQuery.data;
  $: contest = $contestQuery.data;
  $: compClasses = $compClassesQuery.data;
  $: problems = $problemsQuery.data;
  $: ticks = $ticksQuery.data;
  $: selectedCompClass = compClasses?.find(
    ({ id }) => id === contender?.compClassId,
  );
  $: startTime = selectedCompClass?.timeBegin
    ? parseISO(selectedCompClass.timeBegin)
    : new Date(8640000000000000);
  $: endTime = selectedCompClass?.timeEnd
    ? parseISO(selectedCompClass.timeEnd)
    : new Date(-8640000000000000);
  $: gracePeriodEnd = add(endTime, {
    minutes: (contest?.gracePeriod ?? 0) / (1_000_000_000 * 60),
  });

  $: {
    if (startTime && endTime && gracePeriodEnd) {
      calulateAndScheduleNextStateChange();
    }
  }

  $: disabled = ["NOT_STARTED", "ENDED"].includes(state);

  $: {
    if (contender) {
      score = contender.score;
      placement = contender.placement;
    }
  }

  const handleShowTab = ({ detail }: SlTabShowEvent) => {
    if (detail.name === "results") {
      resultsConnected = true;
    }
  };

  const handleVisibilityChange = () => {
    switch (document.visibilityState) {
      case "hidden":
        resultsConnected = false;
        tabGroup.show("problems");

        clearInterval(updateStateTimerId);
        updateStateTimerId = 0;
        break;
      case "visible":
        calulateAndScheduleNextStateChange();
        break;
    }
  };

  const calulateAndScheduleNextStateChange = () => {
    clearInterval(updateStateTimerId);
    updateStateTimerId = 0;

    const now = new Date();
    let durationUntilNextState = 0;

    switch (true) {
      case isBefore(now, startTime):
        state = "NOT_STARTED";
        durationUntilNextState = differenceInMilliseconds(startTime, now);

        break;
      case isBefore(now, endTime):
        state = "RUNNING";
        durationUntilNextState = differenceInMilliseconds(endTime, now);

        break;
      case isBefore(now, gracePeriodEnd):
        state = "GRACE_PERIOD";
        durationUntilNextState = differenceInMilliseconds(gracePeriodEnd, now);

        break;
      default:
        state = "ENDED";
    }

    if (durationUntilNextState) {
      updateStateTimerId = setTimeout(
        calulateAndScheduleNextStateChange,
        durationUntilNextState,
      );
    }
  };

  onMount(() => {
    eventSource = new EventSource(
      `${configData.API_URL}/contests/${$session.contestId}/events`,
    );

    eventSource.addEventListener("CONTENDER_SCORE_UPDATED", (e) => {
      const event = JSON.parse(e.data) as ContenderScoreUpdatedEvent;

      if (event.contenderId === contender?.id) {
        score = event.score;
        placement = event.placement;
      }
    });
  });

  onDestroy(() => {
    eventSource?.close();
    eventSource = undefined;
  });
</script>

<svelte:window on:visibilitychange={handleVisibilityChange} />

{#if !contender || !contest || !compClasses || !problems || !ticks || !selectedCompClass}
  <Loading />
{:else}
  <main>
    <div class="sticky">
      <Header
        registrationCode={$session.registrationCode}
        contestName={contest.name}
        compClassName={selectedCompClass?.name}
        contenderName={contender.name}
        contenderClub={contender.clubName}
        {score}
        {placement}
        {state}
        {startTime}
        {endTime}
        {disabled}
      />
    </div>
    <sl-tab-group bind:this={tabGroup} on:sl-tab-show={handleShowTab}>
      <sl-tab slot="nav" panel="problems">Scorecard</sl-tab>
      <sl-tab slot="nav" panel="results">Results</sl-tab>

      <sl-tab-panel name="problems">
        {#each problems as problem}
          <ProblemView
            {problem}
            tick={ticks.find(({ problemId }) => problemId === problem.id)}
            {disabled}
          />
        {/each}
      </sl-tab-panel>
      <sl-tab-panel name="results">
        {#if resultsConnected && contender.compClassId}
          <ScoreboardProvider contestId={$session.contestId}>
            <ResultList compClassId={contender.compClassId} />
          </ScoreboardProvider>
        {/if}
      </sl-tab-panel>
    </sl-tab-group>
  </main>
{/if}

<style>
  sl-tab-panel::part(base) {
    padding-top: var(--sl-spacing-small);
    padding-bottom: 0;
  }

  main {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .sticky {
    position: sticky;
    top: 0;
    left: 0;
    right: 0;
    z-index: 10;
    background-color: var(--sl-color-primary-200);
    padding: var(--sl-spacing-small);
    padding-bottom: 0;
  }

  sl-tab-group {
    --track-color: transparent;
    padding-inline: var(--sl-spacing-small);
    padding-bottom: var(--sl-spacing-small);
  }

  sl-tab-panel[name="problems"]::part(base) {
    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-x-small);
    width: 100%;
  }
</style>
