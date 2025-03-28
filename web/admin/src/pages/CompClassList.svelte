<script lang="ts">
  import { getCompClassesQuery } from "@climblive/lib/queries";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const compClassesQuery = getCompClassesQuery(contestId);

  let compClasses = $derived($compClassesQuery.data);
</script>

<section>
  {#if compClasses}
    <ul>
      {#each compClasses as compClass (compClass.id)}
        <li>{compClass.name}</li>
      {/each}
    </ul>
  {/if}
</section>

<style>
  section {
    display: flex;
    gap: var(--sl-spacing-x-small);
  }
</style>
