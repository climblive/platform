<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import {
    HoldColorIndicator,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import { type Problem, type ProblemID } from "@climblive/lib/models";
  import {
    getProblemsQuery,
    getTicksByContestQuery,
  } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";
  import DeleteProblem from "./DeleteProblem.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));
  const ticksQuery = $derived(getTicksByContestQuery(contestId));

  const ascentsByProblem = $derived.by(() => {
    const ascentsByProblem = new Map<ProblemID, number>();

    if ($ticksQuery.data === undefined) {
      return ascentsByProblem;
    }

    for (const { problemId } of $ticksQuery.data) {
      const ascents = ascentsByProblem.get(problemId);

      if (ascents !== undefined) {
        ascentsByProblem.set(problemId, ascents + 1);
      } else {
        ascentsByProblem.set(problemId, 1);
      }
    }

    return ascentsByProblem;
  });

  type ProblemWithAscents = Problem & { ascents: number };

  const sortedProblemsWithAscents = $derived.by(() => {
    if ($problemsQuery.data === undefined) {
      return undefined;
    }

    const problems = $problemsQuery.data.map<ProblemWithAscents>((problem) => ({
      ...problem,
      ascents: ascentsByProblem.get(problem.id) ?? 0,
    }));

    return problems?.sort((p1, p2) => p1.number - p2.number);
  });

  const columns: ColumnDefinition<ProblemWithAscents>[] = [
    {
      label: "Number",
      mobile: true,
      render: renderNumberAndColor,
      width: "minmax(max-content, 3fr)",
    },
    {
      label: "Points",
      mobile: true,
      render: renderPoints,
      width: "max-content",
    },
    {
      label: "Flash",
      mobile: true,
      render: renderFlashBonus,
      width: "max-content",
    },
    {
      label: "Tops",
      mobile: false,
      render: renderAscents,
      align: "right",
      width: "max-content",
    },
    {
      mobile: true,
      render: renderControls,
      align: "right",
      width: "max-content",
    },
  ];
</script>

{#snippet renderNumberAndColor({
  number,
  holdColorPrimary,
  holdColorSecondary,
}: Problem)}
  <div class="number">
    <HoldColorIndicator
      --height="1.25rem"
      --width="1.25rem"
      primary={holdColorPrimary}
      secondary={holdColorSecondary}
    />
    № {number}
  </div>
{/snippet}

{#snippet renderPoints({ pointsTop }: ProblemWithAscents)}
  {pointsTop} pts
{/snippet}

{#snippet renderFlashBonus({ flashBonus }: ProblemWithAscents)}
  {#if flashBonus}
    {flashBonus} pts
  {:else}
    -
  {/if}
{/snippet}

{#snippet renderControls({ id }: ProblemWithAscents)}
  <wa-button
    size="small"
    appearance="plain"
    onclick={() => navigate(`/admin/problems/${id}/edit`)}
    label="Edit"
  >
    <wa-icon name="pencil"></wa-icon>
  </wa-button>
  <DeleteProblem problemId={id}>
    {#snippet children({ deleteProblem })}
      <wa-button
        size="small"
        variant="danger"
        appearance="plain"
        onclick={deleteProblem}
        label={`Delete problem ${id}`}
      >
        <wa-icon name="trash"></wa-icon>
      </wa-button>
    {/snippet}
  </DeleteProblem>
{/snippet}

{#snippet renderAscents({ ascents }: ProblemWithAscents)}
  {ascents}
{/snippet}

<section>
  <wa-button
    variant="neutral"
    appearance="accent"
    onclick={() => navigate(`contests/${contestId}/new-problem`)}
    >Create problem</wa-button
  >

  {#if sortedProblemsWithAscents === undefined}
    <Loader />
  {:else if sortedProblemsWithAscents.length > 0}
    <Table {columns} data={sortedProblemsWithAscents} getId={({ id }) => id}
    ></Table>
  {/if}
</section>

<style>
  section {
    display: flex;
    gap: var(--wa-space-xs);
  }

  .number {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }

  section {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-m);
  }
</style>
