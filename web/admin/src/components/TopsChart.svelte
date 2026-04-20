<script lang="ts">
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import "@awesome.me/webawesome/dist/components/scroller/scroller.js";
  import type { ProblemID } from "@climblive/lib/models";
  import {
    getProblemsQuery,
    getTicksByContestQuery,
  } from "@climblive/lib/queries";
  import SubBar from "./SubBar.svelte";

  interface Props {
    contestId: number;
  }

  interface ProblemStats {
    zone1: number;
    zone2: number;
    top: number;
    flash: number;
  }

  const { contestId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));
  const ticksQuery = $derived(getTicksByContestQuery(contestId));

  const statsByProblem = $derived.by(() => {
    const stats = new Map<ProblemID, ProblemStats>();

    if (!ticksQuery.data) {
      return stats;
    }

    for (const tick of ticksQuery.data) {
      let s = stats.get(tick.problemId);

      if (!s) {
        s = { zone1: 0, zone2: 0, top: 0, flash: 0 };
        stats.set(tick.problemId, s);
      }

      if (tick.top && tick.attemptsTop === 1) {
        s.flash++;
      } else if (tick.top) {
        s.top++;
      } else if (tick.zone2) {
        s.zone2++;
      } else if (tick.zone1) {
        s.zone1++;
      }
    }

    return stats;
  });

  const maxCount = $derived.by(() => {
    let max = 0;

    for (const s of statsByProblem.values()) {
      max = Math.max(max, s.zone1, s.zone2, s.top, s.flash);
    }

    return max;
  });

  const sortedProblems = $derived(
    problemsQuery.data
      ? [...problemsQuery.data].sort((a, b) => a.number - b.number)
      : undefined,
  );

  const pct = (count: number) => (maxCount > 0 ? (count / maxCount) * 100 : 0);
</script>

{#if sortedProblems && sortedProblems.length > 0}
  <h3>Tops per problem</h3>
  <wa-scroller orientation="horizontal">
    <div class="chart" role="img" aria-label="Tops per problem">
      {#each sortedProblems as problem (problem.id)}
        {@const stats = statsByProblem.get(problem.id) ?? {
          zone1: 0,
          zone2: 0,
          top: 0,
          flash: 0,
        }}
        {@const z1Pct = pct(stats.zone1)}
        {@const z2Pct = pct(stats.zone2)}
        {@const topPct = pct(stats.top)}
        {@const flashPct = pct(stats.flash)}

        <div class="bar-container">
          <button
            class="bar-stack"
            id="bar-{problem.id}"
            type="button"
            aria-label="Show details for problem #{problem.number}"
          >
            <SubBar
              percentage={z1Pct}
              color={problem.holdColorPrimary}
              opacity={0.35}
            />
            <SubBar
              percentage={z2Pct}
              color={problem.holdColorPrimary}
              opacity={0.55}
            />
            <SubBar
              percentage={topPct}
              color={problem.holdColorPrimary}
              opacity={0.75}
            />
            <SubBar
              percentage={flashPct}
              color={problem.holdColorPrimary}
              opacity={1}
            />
          </button>
          <span class="label">#{problem.number}</span>

          <wa-popover for="bar-{problem.id}" placement="top">
            <div class="popover-body">
              <strong>#{problem.number}</strong>
              {#if problem.zone1Enabled && stats.zone1 + stats.zone2 + stats.top + stats.flash > 0}
                <div>
                  Zone 1: {stats.zone1 + stats.zone2 + stats.top + stats.flash}
                </div>
              {/if}
              {#if problem.zone2Enabled && stats.zone2 + stats.top + stats.flash > 0}
                <div>Zone 2: {stats.zone2 + stats.top + stats.flash}</div>
              {/if}
              <div>Tops: {stats.top + stats.flash}</div>
              {#if stats.flash > 0}
                <div>Flashes: {stats.flash}</div>
              {/if}
            </div>
          </wa-popover>
        </div>
      {/each}
    </div>
  </wa-scroller>
{/if}

<style>
  .chart {
    display: flex;
    gap: var(--wa-space-xs);
    padding-block-start: var(--wa-space-m);
  }

  .bar-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: 0 0 2rem;
  }

  .bar-stack {
    display: flex;
    flex-direction: row;
    align-items: flex-end;
    gap: 2px;
    width: 100%;
    height: 8rem;
    appearance: none;
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
  }

  .label {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    white-space: nowrap;
    margin-block-start: var(--wa-space-xs);
  }

  .popover-body {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-3xs);
    font-size: var(--wa-font-size-s);
  }
</style>
