<script lang="ts">
  import { onMount, type Snippet } from "svelte";

  type T = $$Generic<unknown>;

  export type ColumnDefinition<T> = {
    label?: string;
    mobile: boolean;
    render: (row: T, mobile: boolean) => ReturnType<Snippet>;
    align?: "left" | "right";
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
</script>

<svelte:window onresize={handleResize} />

<table border="0">
  <thead>
    <tr>
      {#each columns as column (column)}
        {#if !mobile || (mobile && column.mobile)}
          <th>{column.label}</th>
        {/if}
      {/each}
    </tr>
  </thead>
  <tbody>
    {#each data as row (getId(row))}
      <tr>
        {#each columns as column, index (index)}
          {#if !mobile || (mobile && column.mobile)}
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
  }

  thead {
    height: 3rem;

    background-color: var(--wa-color-brand-fill-normal);
    color: var(--wa-color-brand-on-normal);

    & th {
      font-weight: var(--wa-font-weight-bold);
      text-align: left;
    }

    & th:first-of-type {
      padding-left: var(--wa-space-m);
    }
  }

  tr {
    height: 3rem;
    cursor: pointer;
  }

  tr:nth-child(even) {
    background-color: var(--wa-color-surface-raised);
  }

  tr:hover {
    background-color: var(--wa-color-surface-raised);
  }

  td:first-of-type {
    padding-left: var(--wa-space-m);
  }

  td[data-align="right"] {
    text-align: right;
  }
</style>
