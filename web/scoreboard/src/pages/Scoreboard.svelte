<script lang="ts">
  import { ResultList, ScoreboardProvider } from "@climblive/lib/components";
  import { getCompClassesQuery, getContestQuery } from "@climblive/lib/queries";
  import { SlSelect } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/option/option.js";
  import "@shoelace-style/shoelace/dist/components/select/select.js";
  import { onMount } from "svelte";
  import Header from "../components/Header.svelte";
  import PoweredBy from "../components/PoweredBy.svelte";
  import Loading from "./Loading.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let compClassSelector: SlSelect | undefined = $state();
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
    {#snippet children({ scoreboard, loading, online })}
      <header data-online={online}>
        <sl-icon name="cloud-slash-fill"></sl-icon>Offline
      </header>
      <main>
        <h1>
          {contest.name}
        </h1>
        <PoweredBy />
        {#if compClasses.length > 1}
          <sl-select
            bind:this={compClassSelector}
            size="small"
            name="compClassId"
            label="Competition class"
            use:value={selectedCompClassId}
            onsl-change={() => {
              selectedCompClassId = Number(compClassSelector?.value);
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
  header {
    background-color: var(--sl-color-danger-600);
    width: 100%;
    height: 2rem;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: var(--sl-spacing-x-small);
    color: white;
    font-weight: var(--sl-font-weight-semibold);
    font-size: var(--sl-font-size-small);

    &[data-online="true"] {
      display: none;
    }
  }

  main {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: var(--sl-spacing-small);
  }

  h1 {
    text-align: center;
    line-height: 1;
    color: black;
    margin-bottom: 0;
  }

  .container {
    margin-top: var(--sl-spacing-small);
    flex-grow: 1;
    display: grid;
    grid-template-columns: repeat(
      var(--num-columns),
      minmax(max-content, 32rem)
    );
    grid-template-rows: 1fr;
    padding-top: 0;
    gap: var(--sl-spacing-small);

    justify-content: center;
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
