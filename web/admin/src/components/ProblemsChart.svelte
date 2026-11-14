<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/scroller/scroller.js";
  import { useEvents } from "@climblive/lib/components";
  import {
    ascentDeregisteredEventSchema,
    ascentRegisteredEventSchema,
    type EventFilter,
    type ProblemID,
    type Tick,
  } from "@climblive/lib/models";
  import {
    getProblemsQuery,
    getTicksByContestQuery,
    removeTickFromQueryCache,
    updateTickInContestQueryCache,
  } from "@climblive/lib/queries";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { onMount } from "svelte";
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

  const queryClient = useQueryClient();
  const { subscribeEvents, unsubscribeEvents } = useEvents();

  const problemsQuery = $derived(getProblemsQuery(contestId));
  const ticksQuery = $derived(getTicksByContestQuery(contestId));

  onMount(() => {
    const currentContestId = contestId;
    const filters: EventFilter[] = [
      {
        contestId: currentContestId,
        eventTypes: {
          ASCENT_REGISTERED: {},
          ASCENT_DEREGISTERED: {},
        },
      },
    ];

    const subscriptionId = subscribeEvents(filters, (eventName, e) => {
      switch (eventName) {
        case "ASCENT_REGISTERED": {
          const event = ascentRegisteredEventSchema.parse(JSON.parse(e.data));
          const tick: Tick = {
            id: event.tickId,
            timestamp: event.timestamp,
            revision: event.revision,
            problemId: event.problemId,
            zone1: event.zone1,
            attemptsZone1: event.attemptsZone1,
            zone2: event.zone2,
            attemptsZone2: event.attemptsZone2,
            top: event.top,
            attemptsTop: event.attemptsTop,
          };

          updateTickInContestQueryCache(queryClient, currentContestId, tick);
          break;
        }
        case "ASCENT_DEREGISTERED": {
          const event = ascentDeregisteredEventSchema.parse(JSON.parse(e.data));

          removeTickFromQueryCache(queryClient, event.tickId);
          break;
        }
      }
    });

    return () => {
      unsubscribeEvents(subscriptionId);
    };
  });

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
      max = Math.max(max, stat.top + stat.flash);
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
    padding-block-end: var(--wa-space-s);
  }

  p {
    font-size: var(--wa-font-size-s);
    color: var(--wa-color-text-quiet);
  }
</style>
