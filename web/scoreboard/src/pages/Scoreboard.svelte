<script lang="ts">
  import logoUrl from "@/static/logo.svg";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import { ResultList, ScoreboardProvider } from "@climblive/lib/components";
  import { value } from "@climblive/lib/forms";
  import { getCompClassesQuery, getContestQuery } from "@climblive/lib/queries";
  import { onMount } from "svelte";
  import Header from "../components/Header.svelte";
  import Loading from "./Loading.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let compClassSelector: WaSelect | undefined = $state();
  let selectedCompClassId: number | undefined = $state();
  let overflow: "pagination" | "scroll" = $state("scroll");

  const contestQuery = $derived(getContestQuery(contestId));
  const compClassesQuery = $derived(getCompClassesQuery(contestId));

  let contest = $derived($contestQuery.data);
  let compClasses = $derived($compClassesQuery.data);

  $effect(() => {
    if (compClasses && !selectedCompClassId) {
      selectedCompClassId = compClasses[0].id;
    }
  });

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
        <wa-icon name="cloud-slash-fill"></wa-icon>Offline
      </header>
      <main>
        <h1>
          {contest.name}
        </h1>
        <p class="logo">
          <img src={logoUrl} alt="ClimbLive" />
        </p>
        {#if compClasses.length > 1}
          <wa-select
            bind:this={compClassSelector}
            size="small"
            name="compClassId"
            label="Competition class"
            {@attach value(selectedCompClassId)}
            onchange={() => {
              selectedCompClassId = Number(compClassSelector?.value);
            }}
          >
            {#each compClasses as compClass (compClass.id)}
              <wa-option value={compClass.id}>{compClass.name}</wa-option>
            {/each}
          </wa-select>
        {/if}
        <div class="container" style="--num-columns: {compClasses.length}">
          {#each compClasses as compClass (compClass.id)}
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
    background-color: var(--wa-color-danger-600);
    width: 100%;
    height: 2rem;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: var(--wa-space-xs);
    color: white;
    font-weight: var(--wa-font-weight-semibold);
    font-size: var(--wa-font-size-s);

    &[data-online="true"] {
      display: none;
    }
  }

  main {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: var(--wa-space-s);
  }

  h1 {
    text-align: center;
    line-height: 1;
    color: var(--wa-color-text-normal);
    margin-bottom: 0;
  }

  .container {
    margin-top: var(--wa-space-s);
    flex-grow: 1;
    display: grid;
    grid-template-columns: repeat(
      var(--num-columns),
      minmax(max-content, 32rem)
    );
    grid-template-rows: 1fr;
    padding-top: 0;
    gap: var(--wa-space-s);

    justify-content: center;
  }

  .class {
    display: flex;
    flex-direction: column;
  }

  wa-select {
    display: none;
  }

  .logo {
    text-align: center;

    & img {
      height: var(--wa-font-size-xl);
    }
  }

  @media screen and (max-width: 512px) {
    h1 {
      font-size: var(--wa-font-size-xl);
    }

    wa-select {
      display: block;
    }

    .container {
      grid-template-columns: 1fr;
    }

    .class[data-selected="false"] {
      display: none;
    }

    .logo > img {
      height: var(--wa-font-size-l);
    }
  }
</style>
