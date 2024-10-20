<script lang="ts">
  import Header from "@/components/Header.svelte";
  import { ResultList, ScoreboardProvider } from "@climblive/lib/components";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import Loading from "./Loading.svelte";

  export let contestId: number;

  const compClassesQuery = getCompClassesQuery(contestId);

  $: compClasses = $compClassesQuery.data;
</script>

{#if !compClasses}
  <Loading />
{:else}
  <ScoreboardProvider {contestId}>
    <main style="grid-template-columns: repeat({compClasses.length}, 1fr)">
      {#each compClasses as compClass}
        <section class="class">
          <Header
            compClassId={compClass.id}
            name={compClass.name}
            timeEnd={new Date(compClass.timeEnd)}
          ></Header>
          <ResultList compClassId={compClass.id} overflow="pagination" />
        </section>
      {/each}
    </main>
  </ScoreboardProvider>
{/if}

<style>
  main {
    display: grid;
    grid-template-rows: 1fr;
    padding: var(--sl-spacing-small);
    gap: var(--sl-spacing-small);
    height: 100%;
  }

  .class {
    display: flex;
    flex-direction: column;
  }
</style>
