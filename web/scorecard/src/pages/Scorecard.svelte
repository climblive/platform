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
  import {
    ascentDeregisteredEventSchema,
    ascentRegisteredEventSchema,
    contenderScoreUpdatedEventSchema,
    type Problem,
    type Tick,
  } from "@climblive/lib/models";
  import {
    getCompClassesQuery,
    getContenderQuery,
    getContestQuery,
    getProblemsQuery,
    getTicksQuery,
    removeTickFromQueryCache,
    updateTickInQueryCache,
  } from "@climblive/lib/queries";
  import { getApiUrl } from "@climblive/lib/utils";
  import type {
    SlRadioGroup,
    SlTabGroup,
    SlTabShowEvent,
  } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/tab-group/tab-group.js";
  import "@shoelace-style/shoelace/dist/components/tab-panel/tab-panel.js";
  import "@shoelace-style/shoelace/dist/components/tab/tab.js";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { add } from "date-fns/add";
  import { getContext, onDestroy, onMount } from "svelte";
  import { type Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const queryClient = useQueryClient();

  const contenderQuery = getContenderQuery($session.contenderId);
  const contestQuery = getContestQuery($session.contestId);
  const compClassesQuery = getCompClassesQuery($session.contestId);
  const problemsQuery = getProblemsQuery($session.contestId);
  const ticksQuery = getTicksQuery($session.contenderId);

  let resultsConnected = $state(false);
  let tabGroup: SlTabGroup | undefined = $state();
  let radioGroup: SlRadioGroup | undefined = $state();
  let eventSource: EventSource | undefined;
  let score: number = $state(0);
  let placement: number | undefined = $state();

  let contender = $derived($contenderQuery.data);
  let contest = $derived($contestQuery.data);
  let compClasses = $derived($compClassesQuery.data);
  let problems = $derived($problemsQuery.data);
  let ticks = $derived($ticksQuery.data);
  let selectedCompClass = $derived(
    compClasses?.find(({ id }) => id === contender?.compClassId),
  );
  let startTime = $derived(
    selectedCompClass?.timeBegin ?? new Date(8640000000000000),
  );
  let endTime = $derived(
    selectedCompClass?.timeEnd ?? new Date(-8640000000000000),
  );
  let gracePeriodEndTime = $derived(
    add(endTime, {
      minutes: (contest?.gracePeriod ?? 0) / (1_000_000_000 * 60),
    }),
  );

  let orderProblemsBy = $state<"number" | "points">("number");

  let sortedProblems = $derived.by<Problem[]>(() => {
    const clonedProblems = [...(problems ?? [])];

    switch (orderProblemsBy) {
      case "number":
        clonedProblems.sort(
          (p1: Problem, p2: Problem) => p1.number - p2.number,
        );

        break;
      case "points":
        clonedProblems.sort(
          (p1: Problem, p2: Problem) =>
            p1.pointsTop +
            (p1.flashBonus ?? 0) -
            p2.pointsTop -
            (p2.flashBonus ?? 0),
        );

        break;
    }

    return clonedProblems;
  });

  let maxProblemNumber = $derived(
    problems?.reduce((max, cur) => {
      return Math.max(max, cur.number);
    }, 0) ?? 0,
  );

  $effect(() => {
    if (contender) {
      score = contender.score?.score ?? 0;
      placement = contender.score?.placement;
    }
  });

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
      `${getApiUrl()}/contenders/${$session.contenderId}/events`,
    );

    eventSource.addEventListener("CONTENDER_SCORE_UPDATED", (e) => {
      const event = contenderScoreUpdatedEventSchema.parse(JSON.parse(e.data));

      if (event.contenderId === contender?.id) {
        score = event.score;
        placement = event.placement;
      }
    });

    eventSource.addEventListener("ASCENT_REGISTERED", (e) => {
      const event = ascentRegisteredEventSchema.parse(JSON.parse(e.data));

      const newTick: Tick = {
        id: event.tickId,
        timestamp: event.timestamp,
        problemId: event.problemId,
        top: event.top,
        attemptsTop: event.attemptsTop,
        zone: event.zone,
        attemptsZone: event.attemptsZone,
      };

      updateTickInQueryCache(queryClient, $session.contenderId, newTick);
    });

    eventSource.addEventListener("ASCENT_DEREGISTERED", (e) => {
      const event = ascentDeregisteredEventSchema.parse(JSON.parse(e.data));

      removeTickFromQueryCache(queryClient, event.tickId);
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

<svelte:window onvisibilitychange={handleVisibilityChange} />

{#if !contender || !contest || !compClasses || !sortedProblems || !ticks || !selectedCompClass}
  <Loading />
{:else}
  <ContestStateProvider {startTime} {endTime} {gracePeriodEndTime}>
    {#snippet children({ contestState })}
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
            {contestState}
            {startTime}
            {endTime}
          />
        </div>
        <sl-tab-group bind:this={tabGroup} onsl-tab-show={handleShowTab}>
          <sl-tab slot="nav" panel="problems">Scorecard</sl-tab>
          <sl-tab slot="nav" panel="results">Results</sl-tab>
          {#if contest.rules}
            <sl-tab slot="nav" panel="info">Info</sl-tab>
          {/if}

          <sl-tab-panel name="problems">
            <sl-radio-group
              size="small"
              bind:this={radioGroup}
              value={orderProblemsBy}
              onsl-change={() => {
                if (radioGroup) {
                  orderProblemsBy = radioGroup.value as typeof orderProblemsBy;
                }
              }}
            >
              <sl-radio-button value="number">
                <sl-icon
                  slot="prefix"
                  name="sort-numeric-down"
                  label="Sort by number"
                ></sl-icon>
                Sort by number
              </sl-radio-button>

              <sl-radio-button value="points" variant="danger">
                <sl-icon
                  slot="prefix"
                  name="sort-down-alt"
                  label="Sort by points"
                ></sl-icon>
                Sort by points
              </sl-radio-button>
            </sl-radio-group>
            {#each sortedProblems as problem}
              <ProblemView
                {problem}
                tick={ticks.find(({ problemId }) => problemId === problem.id)}
                disabled={["NOT_STARTED", "ENDED"].includes(contestState)}
                {maxProblemNumber}
              />
            {/each}
          </sl-tab-panel>
          <sl-tab-panel name="results">
            {#if resultsConnected}
              <ScoreboardProvider contestId={$session.contestId}>
                {#snippet children({ scoreboard, loading })}
                  <ResultList
                    compClassId={selectedCompClass.id}
                    {scoreboard}
                    {loading}
                  />
                {/snippet}
              </ScoreboardProvider>
            {/if}
          </sl-tab-panel>
          <sl-tab-panel name="info">
            <ContestInfo {contest} problems={sortedProblems} {compClasses} />
          </sl-tab-panel>
        </sl-tab-group>
      </main>
    {/snippet}
  </ContestStateProvider>
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
    background-color: #fefefe;
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

  sl-radio-group::part(button-group) {
    width: 100%;
  }

  sl-radio-button {
    flex-grow: 1;

    &::part(button--checked) {
      border-color: var(--sl-color-primary-300);
      background-color: var(--sl-color-primary-200);
      color: var(--sl-color-primary-950);
    }
  }
</style>
