<script lang="ts">
  import Header from "@/components/Header.svelte";
  import { ResultList, ScoreboardProvider } from "@climblive/lib/components";
  import { getCompClassesQuery, getContestQuery } from "@climblive/lib/queries";
  import { add } from "date-fns/add";
  import Loading from "./Loading.svelte";

  export let contestId: number;

  const contestQuery = getContestQuery(contestId);
  const compClassesQuery = getCompClassesQuery(contestId);

  $: contest = $contestQuery.data;
  $: compClasses = $compClassesQuery.data;
</script>

{#if !contest || !compClasses}
  <Loading />
{:else}
  <ScoreboardProvider {contestId}>
    <main>
      <h1>{contest.name}</h1>
      <div
        class="container"
        style="grid-template-columns: repeat({compClasses.length}, 1fr)"
      >
        {#each compClasses as compClass}
          <section class="class">
            <Header
              compClassId={compClass.id}
              name={compClass.name}
              startTime={compClass.timeBegin}
              endTime={compClass.timeEnd}
              gracePeriodEndTime={add(compClass.timeEnd, {
                minutes: contest.gracePeriod / (1_000_000_000 * 60),
              })}
            ></Header>
            <ResultList compClassId={compClass.id} overflow="pagination" />
          </section>
        {/each}
      </div>
    </main>
  </ScoreboardProvider>
{/if}

<style>
  main {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  h1 {
    text-align: center;
    line-height: 1;
  }

  .container {
    flex-grow: 1;
    display: grid;
    grid-template-rows: 1fr;
    padding: var(--sl-spacing-small);
    padding-top: 0;
    gap: var(--sl-spacing-small);
  }

  .class {
    display: flex;
    flex-direction: column;
  }
</style>
