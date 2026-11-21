<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    EmptyState,
    HoldColorIndicator,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import {
    getPointValuesByCompetitorQuery,
    getProblemsQuery,
    getTicksByCompetitorQuery,
  } from "@climblive/lib/queries";
  import { localizedPointsValue } from "@climblive/lib/utils";
  import { format } from "date-fns";

  interface Props {
    competitionId: number;
    competitorId: number;
  }

  const { competitionId, competitorId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(competitionId));
  const ticksQuery = $derived(getTicksByCompetitorQuery(competitorId));
  const pointValuesQuery = $derived(
    getPointValuesByCompetitorQuery(competitorId),
  );

  const pointValues = $derived(
    new Map(
      pointValuesQuery.data?.map(({ problemId, current }) => [
        problemId,
        current,
      ]) ?? [],
    ),
  );

  const problems = $derived(
    new Map(problemsQuery.data?.map((problem) => [problem.id, problem]) ?? []),
  );
  const ticks = $derived(ticksQuery.data);

  type TickAndProblem = { tick: Tick; problem: Problem };

  const tableData = $derived<TickAndProblem[]>(
    ticks
      ?.filter((tick) => tick.zone1 || tick.zone2 || tick.top)
      .map((tick) => {
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

  const columns: ColumnDefinition<TickAndProblem>[] = [
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
    {
      label: "Points",
      mobile: true,
      render: renderPoints,
      width: "max-content",
    },
    {
      label: "Timestamp",
      mobile: false,
      render: renderTimestamp,
      width: "max-content",
    },
  ];
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
