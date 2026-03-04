<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import { type WaSelectEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
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

  let compClasses = $derived(compClassesQuery.data);

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
  <DeleteCompClass compClassId={id}>
    {#snippet children({ deleteCompClass })}
      <wa-dropdown
        onwa-select={(e: WaSelectEvent) => {
          if (e.detail.item.getAttribute("value") === "delete") {
            deleteCompClass();
          } else {
            navigate(`/admin/comp-classes/${id}/edit`);
          }
        }}
      >
        <wa-button slot="trigger" size="small" appearance="plain">
          <wa-icon name="ellipsis-vertical" label="Actions"></wa-icon>
        </wa-button>
        <wa-dropdown-item value="edit">
          <wa-icon slot="icon" name="pencil"></wa-icon>
          Edit
        </wa-dropdown-item>
        <wa-dropdown-item value="delete" variant="danger">
          <wa-icon slot="icon" name="trash"></wa-icon>
          Delete
        </wa-dropdown-item>
      </wa-dropdown>
    {/snippet}
  </DeleteCompClass>
{/snippet}

{#snippet createButton()}
  <wa-button
    variant="neutral"
    appearance="accent"
    onclick={() => navigate(`contests/${contestId}/new-comp-class`)}
    >Create class</wa-button
  >
{/snippet}

<p class="copy">
  Classes represent the categories in which the contenders compete, typically
  divided into Males and Females. The contest duration is defined by the start
  and end times of your classes.
</p>

<section>
  {#if compClasses === undefined}
    <Loader />
  {:else if compClasses.length > 0}
    {@render createButton()}
    <Table
      {columns}
      data={compClasses}
      getId={({ id }) => id}
      onRowClick={({ id }) => navigate(`/admin/comp-classes/${id}/edit`)}
    ></Table>
  {:else}
    <EmptyState
      title="No classes yet"
      description="Create classes to define the categories in which contenders will compete."
    >
      {#snippet actions()}
        {@render createButton()}
      {/snippet}
    </EmptyState>
  {/if}
</section>

<style>
  .copy {
    color: var(--wa-color-text-quiet);
  }

  section {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-m);
  }
</style>
