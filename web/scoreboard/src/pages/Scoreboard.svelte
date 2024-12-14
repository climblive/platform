<script lang="ts">
  import Header from "@/components/Header.svelte";
  import { ResultList, ScoreboardProvider } from "@climblive/lib/components";
  import { getCompClassesQuery, getContestQuery } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/option/option.js";
  import "@shoelace-style/shoelace/dist/components/select/select.js";
  import { onMount } from "svelte";
  import Loading from "./Loading.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let selectedCompClassId: number | undefined = $state();
  let overflow: "pagination" | "scroll" = $state("scroll");

  const contestQuery = getContestQuery(contestId);
  const compClassesQuery = getCompClassesQuery(contestId);

  let contest = $derived($contestQuery.data);
  let compClasses = $derived($compClassesQuery.data);

  $effect(() => {
    if (compClasses && !selectedCompClassId) {
      selectedCompClassId = compClasses[0].id;
    }
  });

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
  onresize={() => {
    determineOverflowBehaviour();
  }}
/>

{#if !contest || !compClasses}
  <Loading />
{:else}
  <ScoreboardProvider {contestId}>
    {#snippet children({ scoreboard, loading })}
      <main>
        <h1>{contest.name}</h1>
        {#if compClasses.length > 1}
          <sl-select
            size="small"
            name="compClassId"
            label="Competition class"
            use:value={selectedCompClassId}
            onsl-change={(event) => {
              selectedCompClassId = parseInt(event.target.value);
            }}
          >
            {#each compClasses as compClass}
              <sl-option value={compClass.id}>{compClass.name}</sl-option>
            {/each}
          </sl-select>
        {/if}
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
                {scoreboard}
              ></Header>
              <ResultList
                compClassId={compClass.id}
                {overflow}
                {scoreboard}
                {loading}
              />
            </section>
          {/each}
        </div>
      </main>
    {/snippet}
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
    color: var(--sl-color-primary-700);
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
    h1 {
      font-size: var(--sl-font-size-x-large);
    }

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
