<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    columns: string[];
    children: Snippet;
  };

  const { columns, children }: Props = $props();
</script>

<table border="0">
  <thead>
    <tr>
      {#each columns as column (column)}
        <th>{column}</th>
      {/each}
    </tr>
  </thead>
  <tbody>
    {@render children()}
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
    border-radius: var(--sl-border-radius-medium);
  }

  thead {
    height: 2rem;

    background-color: var(--sl-color-primary-600);
    color: white;

    & th {
      font-weight: var(--sl-font-weight-normal);
    }

    & th:first-of-type {
      padding-left: var(--sl-spacing-medium);
    }
  }
</style>
