<script lang="ts">
  import ContestInfo from "@/components/ContestInfo.svelte";
  import Header from "@/components/Header.svelte";
  import ProblemView from "@/components/ProblemView.svelte";
  import type { ScorecardSession } from "@/types";
  import {
    ContestStateProvider,
    ResultList,
    ScoreboardProvider,
  } from "@climblive/lib/components";
  import configData from "@climblive/lib/config.json";
  import { contenderScoreUpdatedEventSchema } from "@climblive/lib/models";
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
  import { add } from "date-fns/add";
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
  let tabGroup: SlTabGroup | undefined;
  let eventSource: EventSource | undefined;
  let score: number;
  let placement: number | undefined;

  $: contender = $contenderQuery.data;
  $: contest = $contestQuery.data;
  $: compClasses = $compClassesQuery.data;
  $: problems = $problemsQuery.data;
  $: ticks = $ticksQuery.data;
  $: selectedCompClass = compClasses?.find(
    ({ id }) => id === contender?.compClassId,
  );
  $: startTime = selectedCompClass?.timeBegin ?? new Date(8640000000000000);
  $: endTime = selectedCompClass?.timeEnd ?? new Date(-8640000000000000);
  $: gracePeriodEndTime = add(endTime, {
    minutes: (contest?.gracePeriod ?? 0) / (1_000_000_000 * 60),
  });

  $: {
    if (contender) {
      score = contender.score?.score ?? 0;
      placement = contender.score?.placement;
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
        tearDown();
        break;
      case "visible":
        startEventSubscription();
        break;
    }
  };

  const startEventSubscription = () => {
    if (eventSource) {
      return;
    }

    eventSource = new EventSource(
      `${configData.API_URL}/contenders/${$session.contenderId}/events`,
    );

    eventSource.addEventListener("CONTENDER_SCORE_UPDATED", (e) => {
      const event = contenderScoreUpdatedEventSchema.parse(JSON.parse(e.data));

      if (event.contenderId === contender?.id) {
        score = event.score;
        placement = event.placement;
      }
    });
  };

  const tearDown = () => {
    resultsConnected = false;
    tabGroup?.show("problems");

    eventSource?.close();
    eventSource = undefined;

    stop();
  };

  onMount(() => {
    startEventSubscription();
  });

  onDestroy(() => {
    tearDown();
  });
</script>

<svelte:window on:visibilitychange={handleVisibilityChange} />

{#if !contender || !contest || !compClasses || !problems || !ticks || !selectedCompClass}
  <Loading />
{:else}
  <ContestStateProvider {startTime} {endTime} {gracePeriodEndTime} let:state>
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
        />
      </div>
      <sl-tab-group bind:this={tabGroup} on:sl-tab-show={handleShowTab}>
        <sl-tab slot="nav" panel="problems">Scorecard</sl-tab>
        <sl-tab slot="nav" panel="results">Results</sl-tab>
        {#if contest.rules}
          <sl-tab slot="nav" panel="info">Info</sl-tab>
        {/if}

        <sl-tab-panel name="problems">
          {#each problems as problem}
            <ProblemView
              {problem}
              tick={ticks.find(({ problemId }) => problemId === problem.id)}
              disabled={["NOT_STARTED", "ENDED"].includes(state)}
            />
          {/each}
        </sl-tab-panel>
        <sl-tab-panel name="results">
          {#if resultsConnected && contender.compClassId}
            <ScoreboardProvider
              contestId={$session.contestId}
              let:scoreboard
              let:loading
            >
              <ResultList
                compClassId={contender.compClassId}
                {scoreboard}
                {loading}
              />
            </ScoreboardProvider>
          {/if}
        </sl-tab-panel>
        <sl-tab-panel name="info">
          <ContestInfo {contest} {problems} {compClasses} />
        </sl-tab-panel>
      </sl-tab-group>
    </main></ContestStateProvider
  >
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
