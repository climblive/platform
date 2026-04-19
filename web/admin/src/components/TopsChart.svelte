<script lang="ts">
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import "@awesome.me/webawesome/dist/components/scroller/scroller.js";
  import type { ProblemID } from "@climblive/lib/models";
  import {
    getContendersByContestQuery,
    getProblemsQuery,
    getTicksByContestQuery,
  } from "@climblive/lib/queries";

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
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const totalContenders = $derived(contendersQuery.data?.length ?? 0);

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

  const sortedProblems = $derived(
    problemsQuery.data
      ? [...problemsQuery.data].sort((a, b) => a.number - b.number)
      : undefined,
  );

  const pct = (count: number) =>
    totalContenders > 0 ? (count / totalContenders) * 100 : 0;
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
          <div class="bar-track">
            <div
              class="bar-background"
              style:--bar-color={problem.holdColorPrimary}
            ></div>
            <button
              class="bar-stack"
              id="bar-{problem.id}"
              type="button"
              aria-label="Show details for problem #{problem.number}"
              style:--bar-color={problem.holdColorPrimary}
            >
              {#if flashPct > 0}
                <div
                  class="sub-bar flash"
                  style:--target-height="{flashPct}%"
                ></div>
              {/if}
              {#if topPct > 0}
                <div
                  class="sub-bar top-bar"
                  style:--target-height="{topPct}%"
                ></div>
              {/if}
              {#if z2Pct > 0}
                <div class="sub-bar z2" style:--target-height="{z2Pct}%"></div>
              {/if}
              {#if z1Pct > 0}
                <div class="sub-bar z1" style:--target-height="{z1Pct}%"></div>
              {/if}
            </button>
          </div>
          <span class="label">#{problem.number}</span>

          <wa-popover for="bar-{problem.id}" placement="top">
            <div class="popover-body">
              <strong>#{problem.number}</strong>
              {#if stats.zone1 + stats.zone2 + stats.top + stats.flash > 0}
                <div>
                  Zone 1: {stats.zone1 + stats.zone2 + stats.top + stats.flash}
                </div>
              {/if}
              {#if stats.zone2 + stats.top + stats.flash > 0}
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
    align-items: flex-end;
    gap: var(--wa-space-xs);
    height: 10rem;
    padding-block-start: var(--wa-space-m);
  }

  .bar-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: 0 0 2rem;
    height: 100%;
    justify-content: flex-end;
  }

  .bar-track {
    position: relative;
    width: 100%;
    flex: 1;
  }

  .bar-background {
    position: absolute;
    inset: 0;
    border-radius: var(--wa-border-radius-s) var(--wa-border-radius-s) 0 0;
    background: var(--bar-color);
    opacity: 0.15;
  }

  .bar-stack {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    gap: 2px;
    appearance: none;
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
  }

  .sub-bar {
    width: 100%;
    flex: 0 0 auto;
    background: var(--bar-color);
    animation: grow 0.6s ease-out forwards;
    height: 0;
  }

  .sub-bar:first-child {
    border-radius: var(--wa-border-radius-s) var(--wa-border-radius-s) 0 0;
  }

  .sub-bar.z1 {
    opacity: 0.35;
  }

  .sub-bar.z2 {
    opacity: 0.55;
  }

  .sub-bar.top-bar {
    opacity: 0.75;
  }

  .label {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    white-space: nowrap;
    writing-mode: vertical-rl;
    padding-block-start: var(--wa-space-2xs);
  }

  .popover-body {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-3xs);
    font-size: var(--wa-font-size-s);
  }

  @keyframes grow {
    to {
      height: var(--target-height);
    }
  }
</style>
