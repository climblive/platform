<script lang="ts">
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import Loading from "./Loading.svelte";

  export let contestId: number;

  const compClassesQuery = getCompClassesQuery(contestId);

  $: compClasses = $compClassesQuery.data;
</script>

{#if !compClasses}
  <Loading />
{:else}
  {#each compClasses as compClass}
    <h1>{compClass.name}</h1>
  {/each}
{/if}
