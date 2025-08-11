<script lang="ts">
  import {
    HoldColorIndicator,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Problem } from "@climblive/lib/models";
  import { getProblemsQuery } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";
  import DeleteProblem from "./DeleteProblem.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));

  let problems = $derived($problemsQuery.data);

  const columns: ColumnDefinition<Problem>[] = [
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
      width: "minmax(max-content, 1fr)",
    },
    {
      label: "Flash",
      mobile: true,
      render: renderFlashBonus,
      width: "minmax(max-content, 1fr)",
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

{#snippet renderPoints({ pointsTop }: Problem)}
  {pointsTop} pts
{/snippet}

{#snippet renderFlashBonus({ flashBonus }: Problem)}
  {#if flashBonus}
    {flashBonus} pts
  {:else}
    -
  {/if}
{/snippet}

{#snippet renderControls({ id }: Problem)}
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

<section>
  {#if problems}
    <Table {columns} data={problems} getId={({ id }) => id}></Table>
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
</style>
