<script lang="ts">
  import { onMount, type Snippet } from "svelte";

  type T = $$Generic<unknown>;

  export type ColumnDefinition<T> = {
    label?: string;
    mobile: boolean;
    render: (row: T, mobile: boolean) => ReturnType<Snippet>;
    align?: "left" | "right";
    width?: string;
  };

  type Props<T> = {
    columns: ColumnDefinition<T>[];
    data: T[];
    getId: (row: T) => string | number;
  };

  let mobile = $state(false);

  const { columns, data, getId }: Props<T> = $props();

  function handleResize() {
    mobile = window.innerWidth < 768;
  }

  onMount(() => {
    handleResize();
  });

  const cellVisible = (column: ColumnDefinition<T>) =>
    !mobile || (mobile && column.mobile);

  const gridTemplateColumns = $derived(
    columns
      .map((column) =>
        cellVisible(column) ? (column.width ?? "1fr") : undefined,
      )
      .filter((column) => column !== undefined)
      .join(" "),
  );
</script>

<svelte:window onresize={handleResize} />

<table border="0" style="grid-template-columns: {gridTemplateColumns}">
  <thead>
    <tr>
      {#each columns as column (column)}
        {#if cellVisible(column)}
          <th data-align={column.align ?? "left"}>{column.label}</th>
        {/if}
      {/each}
    </tr>
  </thead>
  <tbody>
    {#each data as row (getId(row))}
      <tr>
        {#each columns as column, index (index)}
          {#if cellVisible(column)}
            <td data-align={column.align ?? "left"}>
              {@render column.render(row, mobile)}
            </td>
          {/if}
        {/each}
      </tr>
    {/each}
  </tbody>
</table>

<style>
  table {
    margin-top: 1rem;
    width: 100%;
    table-layout: fixed;

    border: none;
    border-collapse: separate;
    overflow: hidden;
    border-spacing: 0;
    display: grid;
    grid-template-rows: 1fr;
  }

  thead {
    height: 2.5rem;
    background-color: var(--wa-color-brand-fill-normal);
    color: var(--wa-color-brand-on-normal);

    & th {
      font-weight: var(--wa-font-weight-bold);
      text-align: left;
    }
  }

  tr {
    padding-block: var(--wa-space-xs);
    padding-inline: var(--wa-space-s);
  }

  tbody tr {
    min-height: 3.5rem;
  }

  tbody,
  thead,
  tr {
    display: grid;
    grid-column: 1 / -1;
    grid-template-columns: subgrid;
    column-gap: var(--wa-space-m);
    align-items: center;
  }

  tr:nth-child(even) {
    background-color: var(--wa-color-surface-raised);
  }

  tr:hover {
    background-color: var(--wa-color-surface-raised);
  }

  th[data-align="right"],
  td[data-align="right"] {
    text-align: right;
  }

  td {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
