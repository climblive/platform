<script lang="ts">
  import ContestInfo from "@/components/ContestInfo.svelte";
  import Header from "@/components/Header.svelte";
  import ProblemView from "@/components/ProblemView.svelte";
  import type { ScorecardSession } from "@/types";
  import type { WaTabShowEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/radio-group/radio-group.js";
  import type WaRadioGroup from "@awesome.me/webawesome/dist/components/radio-group/radio-group.js";
  import "@awesome.me/webawesome/dist/components/radio/radio.js";
  import "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
  import type WaTabGroup from "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
  import "@awesome.me/webawesome/dist/components/tab-panel/tab-panel.js";
  import "@awesome.me/webawesome/dist/components/tab/tab.js";
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
  import { useQueryClient } from "@tanstack/svelte-query";
  import { add } from "date-fns/add";
  import { getContext, onDestroy, onMount } from "svelte";
  import { type Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const queryClient = useQueryClient();

  const contenderQuery = $derived(getContenderQuery($session.contenderId));
  const contestQuery = $derived(getContestQuery($session.contestId));
  const compClassesQuery = $derived(getCompClassesQuery($session.contestId));
  const problemsQuery = $derived(getProblemsQuery($session.contestId));
  const ticksQuery = $derived(getTicksQuery($session.contenderId));

  let resultsConnected = $state(false);
  let tabGroup: WaTabGroup | undefined = $state();
  let radioGroup: WaRadioGroup | undefined = $state();
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

  let highestProblemNumber = $derived(
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

  const handleShowTab = ({ detail }: WaTabShowEvent) => {
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
    tabGroup?.setAttribute("active", "problems");

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
        <wa-tab-group bind:this={tabGroup} onwa-tab-show={handleShowTab}>
          <wa-tab slot="nav" panel="problems">Scorecard</wa-tab>
          <wa-tab slot="nav" panel="results">Results</wa-tab>
          <wa-tab slot="nav" panel="info">Info</wa-tab>

          <wa-tab-panel name="problems">
            <wa-radio-group
              orientation="horizontal"
              size="small"
              bind:this={radioGroup}
              value={orderProblemsBy}
              onchange={() => {
                if (radioGroup) {
                  orderProblemsBy = radioGroup.value as typeof orderProblemsBy;
                }
              }}
            >
              <wa-radio value="number" appearance="button">
                <wa-icon name="arrow-down-1-9" label="Sort by number"></wa-icon>
                Sort by number
              </wa-radio>

              <wa-radio value="points" appearance="button">
                <wa-icon name="arrow-up-short-wide" label="Sort by points"
                ></wa-icon>
                Sort by points
              </wa-radio>
            </wa-radio-group>
            {#each sortedProblems as problem (problem.id)}
              <ProblemView
                {problem}
                tick={ticks.find(({ problemId }) => problemId === problem.id)}
                disabled={["NOT_STARTED", "ENDED"].includes(contestState)}
                {highestProblemNumber}
              />
            {/each}
          </wa-tab-panel>
          <wa-tab-panel name="results">
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
          </wa-tab-panel>
          <wa-tab-panel name="info">
            <ContestInfo {contest} problems={sortedProblems} {compClasses} />
          </wa-tab-panel>
        </wa-tab-group>
      </main>
    {/snippet}
  </ContestStateProvider>
{/if}

<style>
  wa-tab-panel::part(base) {
    padding-top: var(--wa-space-s);
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
    background-color: var(--wa-color-surface-default);
    padding: var(--wa-space-s);
  }

  wa-tab-group {
    padding-inline: var(--wa-space-s);
    padding-bottom: var(--wa-space-s);
  }

  wa-tab-panel[name="problems"]::part(base) {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-xs);
  }

  wa-radio {
    flex-grow: 1;
  }

  wa-radio::part(label) {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
  }
</style>
