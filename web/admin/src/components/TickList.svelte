<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    HoldColorIndicator,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import {
    getProblemsQuery,
    getTicksByContenderQuery,
  } from "@climblive/lib/queries";
  import { calculateProblemScore } from "@climblive/lib/utils";
  import { format } from "date-fns";

  interface Props {
    contestId: number;
    contenderId: number;
  }

  const { contestId, contenderId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));
  const ticksQuery = $derived(getTicksByContenderQuery(contenderId));

  let problems = $derived(
    new Map($problemsQuery.data?.map((problem) => [problem.id, problem]) ?? []),
  );
  let ticks = $derived($ticksQuery.data);

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
      .filter((problem) => problem !== undefined)
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
    № {problem.number}
  </div>
{/snippet}

{#snippet renderFlash({ tick }: TickAndProblem)}
  <wa-icon name={tick.top && tick.attemptsTop === 1 ? "bolt" : "minus"}
  ></wa-icon>
{/snippet}

{#snippet renderPoints({ tick, problem }: TickAndProblem)}
  {calculateProblemScore(problem, tick)} pts
{/snippet}

{#snippet renderTimestamp({ tick }: TickAndProblem)}
  {format(tick.timestamp, "yyyy-MM-dd HH:mm")}
{/snippet}

{#if tableData && tableData.length > 0}
  <Table {columns} data={tableData} getId={({ tick }) => tick.id}></Table>
{/if}

<style>
  .number {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }
</style>
