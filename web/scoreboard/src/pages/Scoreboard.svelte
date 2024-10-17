<script lang="ts">
  import {
    ResultList,
    ScoreboardProvider,
    Timer,
  } from "@climblive/lib/components";
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
    <div class="container">
      {#each compClasses as compClass}
        <section class="class">
          <header>
            <h1>{compClass.name}</h1>
            <Timer endTime={new Date(compClass.timeEnd)} />
          </header>
          <ResultList compClassId={compClass.id} />
        </section>
      {/each}
    </div>
  </ScoreboardProvider>
{/if}

<style>
  .container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr;
    padding: var(--sl-spacing-small);
    gap: var(--sl-spacing-small);
  }

  header {
    text-align: center;
    margin-bottom: var(--sl-spacing-large);

    & h1 {
      line-height: var(--sl-line-height-denser);
    }
  }
</style>
