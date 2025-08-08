<script lang="ts">
  import { Table } from "@climblive/lib/components";
  import type { CompClass } from "@climblive/lib/models";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import { format } from "date-fns";
  import type { ColumnDefinition } from "node_modules/@climblive/lib/src/components/Table.svelte";
  import { navigate } from "svelte-routing";
  import DeleteCompClass from "./DeleteCompClass.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const compClassesQuery = $derived(getCompClassesQuery(contestId));

  let compClasses = $derived($compClassesQuery.data);

  const columns: ColumnDefinition<CompClass>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
    },
    {
      label: "Start time",
      mobile: true,
      render: renderTimeBegin,
    },
    {
      label: "End time",
      mobile: false,
      render: renderTimeEnd,
    },
    {
      mobile: true,
      render: renderControls,
      align: "right",
    },
  ];
</script>

{#snippet renderName({ name }: CompClass)}
  {name}
{/snippet}

{#snippet renderTimeBegin({ timeBegin }: CompClass)}
  {format(timeBegin, "yyyy-MM-dd HH:mm")}
{/snippet}

{#snippet renderTimeEnd({ timeEnd }: CompClass)}
  {format(timeEnd, "yyyy-MM-dd HH:mm")}
{/snippet}

{#snippet renderControls({ id }: CompClass)}
  <wa-button
    size="small"
    appearance="plain"
    onclick={() => navigate(`/admin/comp-classes/${id}/edit`)}
    label="Edit"
  >
    <wa-icon name="pencil"></wa-icon>
  </wa-button>
  <DeleteCompClass compClassId={id}>
    {#snippet children({ deleteCompClass })}
      <wa-button
        size="small"
        variant="danger"
        appearance="plain"
        onclick={deleteCompClass}
        label={`Delete comp class ${id}`}
      >
        <wa-icon name="trash"></wa-icon>
      </wa-button>
    {/snippet}
  </DeleteCompClass>
{/snippet}

<section>
  {#if compClasses}
    <Table {columns} data={compClasses} getId={({ id }) => id}></Table>
  {/if}
</section>

<style>
  section {
    display: flex;
    gap: var(--wa-space-xs);
  }
</style>
