<script lang="ts">
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import "@awesome.me/webawesome/dist/components/scroller/scroller.js";
  import type { ProblemID } from "@climblive/lib/models";
  import {
    getProblemsQuery,
    getTicksByContestQuery,
  } from "@climblive/lib/queries";
  import Bar from "./Bar.svelte";

  interface Props {
    contestId: number;
  }

  export interface ProblemStats {
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
      let stat = stats.get(tick.problemId);

      if (!stat) {
        stat = { zone1: 0, zone2: 0, top: 0, flash: 0 };
        stats.set(tick.problemId, stat);
      }

      if (tick.top && tick.attemptsTop === 1) {
        stat.flash++;
      } else if (tick.top) {
        stat.top++;
      } else if (tick.zone2) {
        stat.zone2++;
      } else if (tick.zone1) {
        stat.zone1++;
      }
    }

    return stats;
  });

  const maxCount = $derived.by(() => {
    let max = 0;

    for (const stat of statsByProblem.values()) {
      max = Math.max(max, stat.zone1 + stat.zone2 + stat.top + stat.flash);
    }

    return max;
  });

  const sortedProblems = $derived(
    problemsQuery.data
      ? [...problemsQuery.data].sort((a, b) => a.number - b.number)
      : undefined,
  );
</script>

<section>
  {#if sortedProblems && sortedProblems.length > 0}
    <h2>Tops per problem</h2>
    <p>
      <wa-icon name="lightbulb" variant="regular"></wa-icon> Each problem can be clicked
      to show more details.
    </p>
    <wa-scroller orientation="horizontal">
      {#each sortedProblems as problem (problem.id)}
        {@const stats = statsByProblem.get(problem.id) ?? {
          zone1: 0,
          zone2: 0,
          top: 0,
          flash: 0,
        }}

        <Bar {problem} {stats} {maxCount} />
      {/each}
    </wa-scroller>
  {/if}
</section>

<style>
  wa-scroller::part(content) {
    display: flex;
    gap: var(--wa-space-xs);
  }

  p {
    font-size: var(--wa-font-size-s);
    color: var(--wa-color-text-quiet);
  }
</style>
