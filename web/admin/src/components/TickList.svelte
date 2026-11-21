<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    EmptyState,
    HoldColorIndicator,
    Table,
    useEvents,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import {
    ascentDeregisteredEventSchema,
    ascentRegisteredEventSchema,
    pointValueUpdatedEventSchema,
    type EventFilter,
    type Problem,
    type Tick,
  } from "@climblive/lib/models";
  import {
    getPointValuesByCompetitorQuery,
    getProblemsQuery,
    getTicksByCompetitorQuery,
    removeTickFromQueryCache,
    updatePointValueInQueryCache,
    updateTickInQueryCache,
  } from "@climblive/lib/queries";
  import { localizedPointsValue } from "@climblive/lib/utils";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { format } from "date-fns";
  import { onMount } from "svelte";

  interface Props {
    competitionId: number;
    competitorId: number;
    showPoints: boolean;
  }

  const { competitionId, competitorId, showPoints }: Props = $props();

  const queryClient = useQueryClient();
  const { subscribeEvents, unsubscribeEvents } = useEvents();

  const problemsQuery = $derived(getProblemsQuery(competitionId));
  const ticksQuery = $derived(getTicksByCompetitorQuery(competitorId));
  const pointValuesQuery = $derived(
    showPoints ? getPointValuesByCompetitorQuery(competitorId) : undefined,
  );

  const pointValues = $derived(
    new Map(
      pointValuesQuery?.data?.map(({ problemId, current }) => [
        problemId,
        current,
      ]) ?? [],
    ),
  );

  onMount(() => {
    const currentCompetitorId = competitorId;
    const filters: EventFilter[] = [
      {
        competitionId,
        competitorId: currentCompetitorId,
        eventTypes: {
          ASCENT_REGISTERED: {},
          ASCENT_DEREGISTERED: {},
          POINT_VALUE_UPDATED: {},
        },
      },
    ];

    const subscriptionId = subscribeEvents(filters, (eventName, e) => {
      switch (eventName) {
        case "ASCENT_REGISTERED": {
          const event = ascentRegisteredEventSchema.parse(JSON.parse(e.data));

          if (event.competitorId !== currentCompetitorId) {
            return;
          }

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

          updateTickInQueryCache(queryClient, currentCompetitorId, tick);
          break;
        }
        case "ASCENT_DEREGISTERED": {
          const event = ascentDeregisteredEventSchema.parse(JSON.parse(e.data));

          if (event.competitorId !== currentCompetitorId) {
            return;
          }

          removeTickFromQueryCache(queryClient, event.tickId);
          break;
        }
        case "POINT_VALUE_UPDATED": {
          const event = pointValueUpdatedEventSchema.parse(JSON.parse(e.data));

          if (event.competitorId !== currentCompetitorId) {
            return;
          }

          updatePointValueInQueryCache(queryClient, currentCompetitorId, event);
          break;
        }
      }
    });

    return () => unsubscribeEvents(subscriptionId);
  });

  const problems = $derived(
    new Map(problemsQuery.data?.map((problem) => [problem.id, problem]) ?? []),
  );
  const ticks = $derived(ticksQuery.data);

  type TickAndProblem = { tick: Tick; problem: Problem };

  const tableData = $derived<TickAndProblem[]>(
    ticks
      ?.map((tick) => {
        const problem = problems.get(tick.problemId);
        if (problem === undefined) {
          return undefined;
        }

        return {
          tick,
          problem,
        };
      })
      .filter((entry) => entry !== undefined)
      .sort(
        (t1: TickAndProblem, t2: TickAndProblem) =>
          t1.problem.number - t2.problem.number,
      ) ?? [],
  );

  const columns = $derived<ColumnDefinition<TickAndProblem>[]>([
    {
      label: "Problem",
      mobile: true,
      render: renderProblemNumberAndColor,
      width: "3fr",
    },
    {
      label: "Result",
      mobile: true,
      render: renderResult,
      width: "max-content",
      align: "left",
    },
    {
      label: "Flash",
      mobile: true,
      render: renderFlash,
      width: "max-content",
      align: "left",
    },
    ...(showPoints
      ? [
          {
            label: "Points",
            mobile: true,
            render: renderPoints,
            width: "max-content",
          },
        ]
      : []),
    {
      label: "Timestamp",
      mobile: false,
      render: renderTimestamp,
      width: "max-content",
    },
  ]);
</script>

{#snippet renderProblemNumberAndColor({ problem }: TickAndProblem)}
  <div class="number">
    <HoldColorIndicator
      --height="1.25rem"
      --width="1.25rem"
      primary={problem.holdColorPrimary}
      secondary={problem.holdColorSecondary}
    />
    #{problem.number}
  </div>
{/snippet}

{#snippet renderResult({ tick }: TickAndProblem)}
  {#if tick.top}
    T
  {:else if tick.zone2}
    Z2
  {:else if tick.zone1}
    Z1
  {:else}
    <wa-icon name="minus"></wa-icon>
  {/if}
{/snippet}

{#snippet renderFlash({ tick }: TickAndProblem)}
  <wa-icon name={tick.top && tick.attemptsTop === 1 ? "bolt" : "minus"}
  ></wa-icon>
{/snippet}

{#snippet renderTimestamp({ tick }: TickAndProblem)}
  {format(tick.timestamp, "yyyy-MM-dd HH:mm")}
{/snippet}

{#snippet renderPoints({ tick }: TickAndProblem)}
  {@const points = pointValues.get(tick.problemId)}
  {#if points !== undefined}
    {localizedPointsValue(points)}
  {:else}
    —
  {/if}
{/snippet}

{#if tableData && tableData.length > 0}
  <Table {columns} data={tableData} getId={({ tick }) => tick.id}></Table>
{:else if tableData}
  <EmptyState
    title="No results yet"
    description="Problems will appear here once the competitor starts logging attempts."
  />
{/if}

<style>
  wa-icon[name="bolt"] {
    color: var(--wa-color-yellow-50);
  }

  .number {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }
</style>
