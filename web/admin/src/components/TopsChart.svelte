<script lang="ts">
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

  const { contestId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));
  const ticksQuery = $derived(getTicksByContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const totalContenders = $derived(contendersQuery.data?.length ?? 0);

  const topsByProblem = $derived.by(() => {
    const tops = new Map<ProblemID, number>();

    if (!ticksQuery.data) {
      return tops;
    }

    for (const tick of ticksQuery.data) {
      if (tick.top) {
        tops.set(tick.problemId, (tops.get(tick.problemId) ?? 0) + 1);
      }
    }

    return tops;
  });

  const sortedProblems = $derived(
    problemsQuery.data
      ? [...problemsQuery.data].sort((a, b) => a.number - b.number)
      : undefined,
  );
</script>

{#if sortedProblems && sortedProblems.length > 0}
  <wa-scroller orientation="horizontal">
    <div class="chart" role="img" aria-label="Tops per problem">
      {#each sortedProblems as problem (problem.id)}
        {@const tops = topsByProblem.get(problem.id) ?? 0}
        {@const heightPercent =
          totalContenders > 0 ? (tops / totalContenders) * 100 : 0}

        <div class="bar-container">
          <span class="count">{tops}</span>
          <div class="bar-track">
            <div
              class="bar-background"
              style:--bar-color={problem.holdColorPrimary}
            ></div>
            <div
              class="bar"
              style:--target-height="{heightPercent}%"
              style:--bar-color={problem.holdColorPrimary}
              style:--bar-color-secondary={problem.holdColorSecondary ||
                problem.holdColorPrimary}
            ></div>
          </div>
          <span class="label">#{problem.number}</span>
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
    gap: 0.25rem;
  }

  .count {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
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

  .bar {
    position: absolute;
    bottom: 0;
    width: 100%;
    border-radius: var(--wa-border-radius-s) var(--wa-border-radius-s) 0 0;
    background: linear-gradient(
      to bottom,
      var(--bar-color) 50%,
      var(--bar-color-secondary) 100%
    );
    animation: grow 0.6s ease-out forwards;
    height: 0;
  }

  .label {
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    white-space: nowrap;
  }

  @keyframes grow {
    to {
      height: var(--target-height);
    }
  }
</style>
