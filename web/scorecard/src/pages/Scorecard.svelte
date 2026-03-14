<script lang="ts">
  import ContestInfo from "@/components/ContestInfo.svelte";
  import Header from "@/components/Header.svelte";
  import ProblemView from "@/components/ProblemView.svelte";
  import Summary from "@/components/Summary.svelte";
  import SummaryCards from "@/components/SummaryCards.svelte";
  import type { ScorecardSession } from "@/types";
  import type { WaTabShowEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
  import type WaTabGroup from "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
  import "@awesome.me/webawesome/dist/components/tab-panel/tab-panel.js";
  import "@awesome.me/webawesome/dist/components/tab/tab.js";
  import {
      ContestStateProvider,
      EmptyState,
      ResultList,
      ScoreboardProvider,
      SplashScreen,
  } from "@climblive/lib/components";
  import {
      ascentDeregisteredEventSchema,
      ascentRegisteredEventSchema,
      contenderPublicInfoUpdatedEventSchema,
      contenderScoreUpdatedEventSchema,
      raffleWinnerDrawnEventSchema,
      type Problem,
      type Tick,
  } from "@climblive/lib/models";
  import {
      getCompClassesQuery,
      getContenderQuery,
      getContestQuery,
      getProblemsQuery,
      getTicksByContenderQuery,
      removeTickFromQueryCache,
      updateContenderPublicInfoInQueryCache,
      updateTickInQueryCache,
  } from "@climblive/lib/queries";
  import { getApiUrl } from "@climblive/lib/utils";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { add } from "date-fns/add";
  import { getContext, onDestroy, onMount } from "svelte";
  import { type Readable } from "svelte/store";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const queryClient = useQueryClient();

  const contenderQuery = $derived(getContenderQuery($session.contenderId));
  const contestQuery = $derived(getContestQuery($session.contestId));
  const compClassesQuery = $derived(getCompClassesQuery($session.contestId));
  const problemsQuery = $derived(getProblemsQuery($session.contestId));
  const ticksQuery = $derived(getTicksByContenderQuery($session.contenderId));

  let resultsConnected = $state(false);
  let tabGroup: WaTabGroup | undefined = $state();
  let raffleWinnerDialog: WaDialog | undefined = $state();
  let stickyHeader: HTMLDivElement | undefined = $state();
  let eventSource: EventSource | undefined;
  let score: number = $state(0);
  let placement: number | undefined = $state();
  let finalist: boolean = $state(false);

  let contender = $derived(contenderQuery.data);
  let contest = $derived(contestQuery.data);
  let compClasses = $derived(compClassesQuery.data);
  let problems = $derived(problemsQuery.data);
  let ticks = $derived(ticksQuery.data);
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
  let sortDirection = $state<"asc" | "desc">("asc");

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

    if (sortDirection === "desc") {
      clonedProblems.reverse();
    }

    return clonedProblems;
  });

  let numberSortIcon = $derived(
    orderProblemsBy === "number" && sortDirection === "desc"
      ? "arrow-up-9-1"
      : "arrow-down-1-9",
  );

  let numberSortLabel = $derived(
    orderProblemsBy === "number" && sortDirection === "desc"
      ? "Sort by number descending"
      : "Sort by number ascending",
  );

  let pointsSortIcon = $derived(
    orderProblemsBy === "points" && sortDirection === "desc"
      ? "arrow-down-wide-short"
      : "arrow-up-short-wide",
  );

  let pointsSortLabel = $derived(
    orderProblemsBy === "points" && sortDirection === "desc"
      ? "Sort by points descending"
      : "Sort by points ascending",
  );

  let highestProblemNumber = $derived(
    problems?.reduce((max, cur) => {
      return Math.max(max, cur.number);
    }, 0) ?? 0,
  );

  $effect(() => {
    if (contender) {
      score = contender.score?.score ?? 0;
      placement = contender.score?.placement;
      finalist = contender.score?.finalist ?? false;
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

    eventSource.addEventListener("CONTENDER_PUBLIC_INFO_UPDATED", (e) => {
      const event = contenderPublicInfoUpdatedEventSchema.parse(
        JSON.parse(e.data),
      );

      if (event.contenderId !== contender?.id) {
        return;
      }

      updateContenderPublicInfoInQueryCache(queryClient, event.contenderId, {
        compClassId: event.compClassId,
        name: event.name,
        withdrawnFromFinals: event.withdrawnFromFinals,
        disqualified: event.disqualified,
      });
    });

    eventSource.addEventListener("CONTENDER_SCORE_UPDATED", (e) => {
      const event = contenderScoreUpdatedEventSchema.parse(JSON.parse(e.data));

      if (event.contenderId === contender?.id) {
        score = event.score;
        placement = event.placement;
        finalist = event.finalist;
      }
    });

    eventSource.addEventListener("ASCENT_REGISTERED", (e) => {
      const event = ascentRegisteredEventSchema.parse(JSON.parse(e.data));

      const newTick: Tick = {
        id: event.tickId,
        timestamp: event.timestamp,
        problemId: event.problemId,
        zone1: event.zone1,
        attemptsZone1: event.attemptsZone1,
        zone2: event.zone2,
        attemptsZone2: event.attemptsZone2,
        top: event.top,
        attemptsTop: event.attemptsTop,
      };

      updateTickInQueryCache(queryClient, $session.contenderId, newTick);
    });

    eventSource.addEventListener("ASCENT_DEREGISTERED", (e) => {
      const event = ascentDeregisteredEventSchema.parse(JSON.parse(e.data));

      removeTickFromQueryCache(queryClient, event.tickId);
    });

    eventSource.addEventListener("RAFFLE_WINNER_DRAWN", (e) => {
      const event = raffleWinnerDrawnEventSchema.parse(JSON.parse(e.data));

      if (event.contenderId === contender?.id && raffleWinnerDialog) {
        raffleWinnerDialog.open = true;
      }
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

  $effect(() => {
    if (!stickyHeader) return;

    const observer = new ResizeObserver(([entry]) => {
      const height = entry.borderBoxSize[0].blockSize;
      stickyHeader!.parentElement?.style.setProperty(
        "--header-height",
        `${height}px`,
      );
    });

    observer.observe(stickyHeader);

    return () => observer.disconnect();
  });

  let showSplash = $state(true);
</script>

<svelte:window onvisibilitychange={handleVisibilityChange} />

{#if showSplash || !contender || !contest || !compClasses || !sortedProblems || !ticks || !selectedCompClass}
  <SplashScreen onComplete={() => (showSplash = false)} />
{:else}
  <ContestStateProvider {startTime} {endTime} {gracePeriodEndTime}>
    {#snippet children({ contestState })}
      <main>
        <div class="sticky" bind:this={stickyHeader}>
          <Header
            registrationCode={$session.registrationCode}
            contestName={contest.name}
            compClassName={selectedCompClass?.name}
            contenderName={contender.name}
            {contestState}
          />
        </div>
        <wa-tab-group bind:this={tabGroup} onwa-tab-show={handleShowTab}>
          <wa-tab slot="nav" panel="problems">Scorecard</wa-tab>
          <wa-tab slot="nav" panel="results">Results</wa-tab>
          <wa-tab slot="nav" panel="info">Info</wa-tab>

          <wa-tab-panel name="problems">
            <SummaryCards
              {score}
              {placement}
              disqualified={contender.disqualified}
              {contestState}
              {startTime}
              {endTime}
              {ticks}
              totalProblems={sortedProblems.length}
            />
            <div class="problems-header">
              <span class="problems-label">Problems</span>
              <div class="sort-buttons">
                <button
                  class="sort-btn"
                  class:active={orderProblemsBy === "number"}
                  onclick={() => {
                    if (orderProblemsBy === "number") {
                      sortDirection = sortDirection === "asc" ? "desc" : "asc";
                    } else {
                      sortDirection = "asc";
                      orderProblemsBy = "number";
                    }
                  }}
                  disabled={problems === undefined || problems.length === 0}
                >
                  <wa-icon name={numberSortIcon} label={numberSortLabel}
                  ></wa-icon>
                  #
                </button>
                <button
                  class="sort-btn"
                  class:active={orderProblemsBy === "points"}
                  onclick={() => {
                    if (orderProblemsBy === "points") {
                      sortDirection = sortDirection === "asc" ? "desc" : "asc";
                    } else {
                      sortDirection = "asc";
                      orderProblemsBy = "points";
                    }
                  }}
                  disabled={problems === undefined || problems.length === 0}
                >
                  <wa-icon name={pointsSortIcon} label={pointsSortLabel}
                  ></wa-icon>
                  Pts
                </button>
              </div>
            </div>
            {#if sortedProblems.length === 0}
              <EmptyState
                title="No problems"
                description="The organizer has not added any problems to this contest yet."
              />
            {:else}
              {#each sortedProblems as problem (problem.id)}
                <ProblemView
                  {problem}
                  tick={ticks.find(({ problemId }) => problemId === problem.id)}
                  disabled={["NOT_STARTED", "ENDED"].includes(contestState)}
                  {highestProblemNumber}
                  disqualified={contender.disqualified}
                />
              {/each}
            {/if}
          </wa-tab-panel>
          <wa-tab-panel name="results">
            {#if contestState !== "NOT_STARTED"}
              <Summary
                {ticks}
                problems={sortedProblems}
                {score}
                {placement}
                {finalist}
                disqualified={contender.disqualified}
              />
            {/if}
            {#if resultsConnected}
              <ScoreboardProvider
                contestId={$session.contestId}
                hideDisqualified
              >
                {#snippet children({ scoreboard, loading })}
                  <ResultList
                    compClassId={selectedCompClass.id}
                    {scoreboard}
                    {loading}
                    highlightedContenderId={contender.id}
                    autoScroll={false}
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

<wa-dialog
  bind:this={raffleWinnerDialog}
  without-header
  class="raffle-winner-dialog"
>
  <h2>Congratulations!</h2>
  <p>You just won a prize in a raffle!</p>
  <wa-button
    slot="footer"
    variant="success"
    appearance="accent"
    size="small"
    onclick={() => {
      if (raffleWinnerDialog) {
        raffleWinnerDialog.open = false;
      }
    }}
  >
    Awesome!
    <wa-icon slot="start" name="gift"></wa-icon>
  </wa-button>
</wa-dialog>

<style>
  wa-tab-panel::part(base) {
    padding-top: var(--wa-space-m);
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
    z-index: 11;
    background-color: var(--wa-color-surface-default);
    padding: 0 var(--wa-space-m);
  }

  wa-tab-group {
    --track-width: 2px;
    padding-inline: var(--wa-space-m);
    padding-bottom: var(--wa-space-m);
  }

  wa-tab-group::part(nav) {
    position: sticky;
    top: var(--header-height, 0px);
    z-index: 10;
    background-color: var(--wa-color-surface-default);
  }

  wa-tab-panel[name="problems"]::part(base) {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-xs);
  }

  .sort-buttons {
    display: flex;
    gap: var(--wa-space-3xs);
  }

  .sort-btn {
    display: inline-flex;
    align-items: center;
    gap: var(--wa-space-3xs);
    font-family: inherit;
    font-size: var(--wa-font-size-2xs);
    font-weight: var(--wa-font-weight-semibold);
    line-height: 1;
    padding: var(--wa-space-3xs) var(--wa-space-2xs);
    border: var(--wa-border-width-s) solid var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-s);
    background: transparent;
    color: var(--wa-color-text-quiet);
    cursor: pointer;
  }

  .sort-btn.active {
    background: var(--wa-color-neutral-fill-quiet-hover);
    color: var(--wa-color-text-default);
  }

  .sort-btn:disabled {
    opacity: 0.4;
    cursor: default;
  }

  .problems-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--wa-space-xs);
  }

  .problems-label {
    font-size: var(--wa-font-size-xs);
    font-weight: var(--wa-font-weight-bold);
    text-transform: uppercase;
    color: var(--wa-color-text-quiet);
    letter-spacing: 0.05em;
  }

  .raffle-winner-dialog {
    &::part(body) {
      text-align: center;
    }

    & wa-button {
      width: 100%;
    }
  }
</style>
