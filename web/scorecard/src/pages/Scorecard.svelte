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
    type Tick,
  } from "@climblive/lib/models";
  import {
    getCompClassesQuery,
    getContenderQuery,
    getContestQuery,
    getProblemsQuery,
    getTicksQuery,
  } from "@climblive/lib/queries";
  import { getApiUrl } from "@climblive/lib/utils";
  import type { SlTabGroup, SlTabShowEvent } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/tab-group/tab-group.js";
  import "@shoelace-style/shoelace/dist/components/tab-panel/tab-panel.js";
  import "@shoelace-style/shoelace/dist/components/tab/tab.js";
  import { useQueryClient, type QueryKey } from "@tanstack/svelte-query";
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

      const queryKey: QueryKey = [
        "ticks",
        { contenderId: $session.contenderId },
      ];

      queryClient.setQueryData<Tick[]>(queryKey, (oldTicks) => {
        const newTick: Tick = {
          id: event.tickId,
          timestamp: event.timestamp,
          problemId: event.problemId,
          top: event.top,
          attemptsTop: event.attemptsTop,
          zone: event.zone,
          attemptsZone: event.attemptsZone,
        };

        const predicate = (tick: Tick) => tick.problemId === event.problemId;

        const found = (oldTicks ?? []).findIndex(predicate) !== -1;

        if (found) {
          return (oldTicks ?? []).map((oldTick) => {
            if (predicate(oldTick)) {
              return newTick;
            } else {
              return oldTick;
            }
          });
        } else {
          return [...(oldTicks ?? []), newTick];
        }
      });
    });

    eventSource.addEventListener("ASCENT_DEREGISTERED", (e) => {
      const event = ascentDeregisteredEventSchema.parse(JSON.parse(e.data));

      const queryKey = ["ticks"];
      queryClient.setQueriesData<Tick[]>(
        {
          queryKey,
          exact: false,
        },
        (oldTicks) => {
          const predicate = (tick: Tick) => tick.id !== event.tickId;

          return oldTicks ? oldTicks.filter(predicate) : undefined;
        },
      );
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

{#if !contender || !contest || !compClasses || !problems || !ticks || !selectedCompClass}
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
            {#each problems as problem}
              <ProblemView
                {problem}
                tick={ticks.find(({ problemId }) => problemId === problem.id)}
                disabled={["NOT_STARTED", "ENDED"].includes(contestState)}
              />
            {/each}
          </sl-tab-panel>
          <sl-tab-panel name="results">
            {#if resultsConnected && contender.compClassId}
              <ScoreboardProvider contestId={$session.contestId}>
                {#snippet children({ scoreboard, loading })}
                  <ResultList
                    compClassId={contender.compClassId}
                    {scoreboard}
                    {loading}
                  />
                {/snippet}
              </ScoreboardProvider>
            {/if}
          </sl-tab-panel>
          <sl-tab-panel name="info">
            <ContestInfo {contest} {problems} {compClasses} />
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
