<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import Loader from "@/components/Loader.svelte";
  import { type WaSelectEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import WaDropdownItem from "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import "@awesome.me/webawesome/dist/components/tooltip/tooltip.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Category } from "@climblive/lib/models";
  import { getCategoriesQuery } from "@climblive/lib/queries";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";
  import DeleteCategory from "./DeleteCategory.svelte";

  const maxCategories = 20;
  const createButtonId = $props.id();

  interface Props {
    competitionId: string;
  }

  const { competitionId }: Props = $props();

  const categoriesQuery = $derived(getCategoriesQuery(competitionId));

  let categories = $derived(categoriesQuery.data);

  const limitReached = $derived(
    categories !== undefined && categories.length >= maxCategories,
  );

  const columns: ColumnDefinition<Category>[] = [
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

{#snippet renderName({ name }: Category)}
  {name}
{/snippet}

{#snippet renderDescription({ description }: Category)}
  {description}
{/snippet}

{#snippet renderTimeBegin({ timeBegin }: Category)}
  {format(timeBegin, "yyyy-MM-dd HH:mm")}
{/snippet}

{#snippet renderTimeEnd({ timeEnd }: Category)}
  {format(timeEnd, "yyyy-MM-dd HH:mm")}
{/snippet}

{#snippet renderControls({ id }: Category)}
  <DeleteCategory categoryId={id}>
    {#snippet children({ deleteCategory })}
      <wa-dropdown
        onwa-select={(event: WaSelectEvent) => {
          if ((event.detail.item as WaDropdownItem).value === "delete") {
            deleteCategory();
          } else {
            navigate(`/admin/categories/${id}/edit`);
          }
        }}
      >
        <wa-button slot="trigger" size="s" appearance="plain">
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
  </DeleteCategory>
{/snippet}

{#snippet createButton()}
  <wa-button
    id={createButtonId}
    variant="neutral"
    appearance="accent"
    disabled={limitReached}
    onclick={() => navigate(`competitions/${competitionId}/new-category`)}
    >Create category</wa-button
  >
  {#if limitReached}
    <wa-tooltip for={createButtonId} placement="top-start"
      >Limit of {maxCategories} categories per competition reached.</wa-tooltip
    >
  {/if}
{/snippet}

<p class="copy">
  Categories represent the classes in which the competitors compete, typically
  divided into Males and Females. The competition duration is defined by the
  start and end times of your categories.
</p>

<section>
  {#if categories === undefined}
    <Loader />
  {:else if categories.length > 0}
    {@render createButton()}
    <Table {columns} data={categories} getId={({ id }) => id}></Table>
  {:else}
    <EmptyState
      title="No categories yet"
      description="Add categories for competitors to compete in."
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
