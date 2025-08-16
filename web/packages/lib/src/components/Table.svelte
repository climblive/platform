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
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-neutral-border-quiet);
    border-radius: var(--wa-border-radius-m);
  }

  @supports (grid-template-columns: subgrid) {
    table {
      display: grid;
      grid-template-rows: 1fr;
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
  }

  th:first-of-type,
  td:first-of-type {
    padding-inline-start: var(--wa-space-s);
  }

  th:last-of-type,
  td:last-of-type {
    padding-inline-end: var(--wa-space-s);
  }

  thead {
    height: 3rem;
    background-color: var(--wa-color-neutral-fill-quiet);
    color: var(--wa-color-neutral-on-normal);

    & th {
      font-weight: var(--wa-font-weight-bold);
      text-align: left;
    }
  }

  tbody tr {
    height: 3.5rem;
  }

  tbody tr:not(:last-of-type) {
    height: 3.5rem;
    border-bottom: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-raised);
  }

  tbody tr:hover {
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
