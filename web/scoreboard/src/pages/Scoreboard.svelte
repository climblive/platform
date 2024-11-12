<script lang="ts">
  import Header from "@/components/Header.svelte";
  import { ResultList, ScoreboardProvider } from "@climblive/lib/components";
  import { getCompClassesQuery, getContestQuery } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/option/option.js";
  import "@shoelace-style/shoelace/dist/components/select/select.js";
  import { onMount } from "svelte";
  import Loading from "./Loading.svelte";

  export let contestId: number;

  let selectedCompClassId: number | undefined;
  let overflow: "pagination" | "scroll" = "scroll";

  const contestQuery = getContestQuery(contestId);
  const compClassesQuery = getCompClassesQuery(contestId);

  $: contest = $contestQuery.data;
  $: compClasses = $compClassesQuery.data;

  $: {
    if (compClasses && !selectedCompClassId) {
      selectedCompClassId = compClasses[0].id;
    }
  }

  const value = (node: HTMLElement, value: string | number | undefined) => {
    node.setAttribute("value", value?.toString() ?? "");

    return {
      update(value: string | number | undefined) {
        node.setAttribute("value", value?.toString() ?? "");
      },
    };
  };

  const determineOverflowBehaviour = () => {
    overflow = window.innerWidth <= 512 ? "scroll" : "pagination";
  };

  onMount(() => {
    determineOverflowBehaviour();
  });
</script>

<svelte:window
  on:resize={() => {
    determineOverflowBehaviour();
  }}
/>

{#if !contest || !compClasses}
  <Loading />
{:else}
  <ScoreboardProvider {contestId}>
    <main>
      <h1>{contest.name}</h1>
      <sl-select
        size="small"
        name="compClassId"
        label="Competition class"
        use:value={selectedCompClassId}
        on:sl-change={(event) => {
          selectedCompClassId = parseInt(event.target.value);
        }}
      >
        {#each compClasses as compClass}
          <sl-option value={compClass.id}>{compClass.name}</sl-option>
        {/each}
      </sl-select>
      <div class="container" style="--num-columns: {compClasses.length}">
        {#each compClasses as compClass}
          <section
            class="class"
            data-selected={compClass.id === selectedCompClassId
              ? "true"
              : "false"}
          >
            <Header
              compClassId={compClass.id}
              name={compClass.name}
              startTime={compClass.timeBegin}
              endTime={compClass.timeEnd}
            ></Header>
            <ResultList compClassId={compClass.id} {overflow} />
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
    padding: var(--sl-spacing-small);
    gap: var(--sl-spacing-small);
  }

  h1 {
    text-align: center;
    line-height: 1;
  }

  .container {
    flex-grow: 1;
    display: grid;
    grid-template-columns: repeat(var(--num-columns), 1fr);
    grid-template-rows: 1fr;
    padding-top: 0;
    gap: var(--sl-spacing-small);
  }

  .class {
    display: flex;
    flex-direction: column;
  }

  sl-select {
    display: none;
  }

  @media screen and (max-width: 512px) {
    sl-select {
      display: block;
    }

    .container {
      grid-template-columns: 1fr;
    }

    .class[data-selected="false"] {
      display: none;
    }
  }
</style>
