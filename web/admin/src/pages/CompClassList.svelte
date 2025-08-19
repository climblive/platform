<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { CompClass } from "@climblive/lib/models";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";
  import DeleteCompClass from "./DeleteCompClass.svelte";

  interface Props {
    contestId: number;
  }

  const { contestId }: Props = $props();

  const compClassesQuery = $derived(getCompClassesQuery(contestId));

  let compClasses = $derived($compClassesQuery.data);

  const columns: ColumnDefinition<CompClass>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "1fr",
    },
    {
      label: "Description",
      mobile: false,
      render: renderDescription,
      width: "1fr",
    },
    {
      label: "Start time",
      mobile: true,
      render: renderTimeBegin,
      width: "max-content",
    },
    {
      label: "End time",
      mobile: false,
      render: renderTimeEnd,
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

{#snippet renderName({ name }: CompClass)}
  {name}
{/snippet}

{#snippet renderDescription({ description }: CompClass)}
  {description}
{/snippet}

{#snippet renderTimeBegin({ timeBegin }: CompClass)}
  {format(timeBegin, "yyyy-MM-dd HH:mm")}
{/snippet}

{#snippet renderTimeEnd({ timeEnd }: CompClass)}
  {format(timeEnd, "yyyy-MM-dd HH:mm")}
{/snippet}

{#snippet renderControls({ id }: CompClass)}
  <div class="controls">
    <wa-button
      size="small"
      appearance="plain"
      onclick={() => navigate(`/admin/comp-classes/${id}/edit`)}
    >
      <wa-icon name="pencil" label="Edit"></wa-icon>
    </wa-button>
    <DeleteCompClass compClassId={id}>
      {#snippet children({ deleteCompClass })}
        <wa-button
          size="small"
          variant="danger"
          appearance="plain"
          onclick={deleteCompClass}
        >
          <wa-icon name="trash" label={`Delete comp class ${id}`}></wa-icon>
        </wa-button>
      {/snippet}
    </DeleteCompClass>
  </div>
{/snippet}

<section>
  <p class="copy">
    Classes represent the categories in which the contenders compete, typically
    divided into Males and Females. The contest duration is defined by the start
    and end times of your classes.
  </p>

  <wa-button
    variant="brand"
    appearance="accent"
    onclick={() => navigate(`contests/${contestId}/new-comp-class`)}
    >Create</wa-button
  >

  {#if compClasses?.length}
    <Table {columns} data={compClasses} getId={({ id }) => id}></Table>
  {/if}
</section>

<style>
  .controls {
    display: flex;
  }

  .copy {
    color: var(--wa-color-text-quiet);
  }
</style>
