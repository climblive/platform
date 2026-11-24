<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import {
    FullLogo,
    ResultList,
    ScoreboardProvider,
    SplashScreen,
  } from "@climblive/lib/components";
  import { value } from "@climblive/lib/forms";
  import { formatCategoryEligibility } from "@climblive/lib/utils";
  import {
    getCategoriesQuery,
    getCompetitionQuery,
  } from "@climblive/lib/queries";
  import { onMount } from "svelte";
  import Header from "../components/Header.svelte";

  interface Props {
    competitionId: number;
  }

  let { competitionId }: Props = $props();

  let categorySelector: WaSelect | undefined = $state();
  let selectedCategoryId: number | undefined = $state();
  let overflow: "pagination" | "scroll" = $state("scroll");

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const categoriesQuery = $derived(getCategoriesQuery(competitionId));

  let competition = $derived(competitionQuery.data);
  let categories = $derived(categoriesQuery.data);

  $effect(() => {
    if (categories && categories.length > 0 && !selectedCategoryId) {
      selectedCategoryId = categories[0].id;
    }
  });

  const determineOverflowBehaviour = () => {
    overflow = window.innerWidth <= 512 ? "scroll" : "pagination";
  };

  onMount(() => {
    determineOverflowBehaviour();
  });

  let showSplash = $state(true);
</script>

<svelte:window
  onresize={() => {
    determineOverflowBehaviour();
  }}
/>

{#if showSplash || !competition || !categories}
  <SplashScreen onComplete={() => (showSplash = false)} />
{:else}
  <ScoreboardProvider {competitionId} hideDisqualified>
    {#snippet children({ scoreboard, loading, online })}
      <header data-online={online}>
        <wa-icon name="triangle-exclamation"></wa-icon>Offline
      </header>
      <main>
        <div class="banner">
          <h1>
            {competition.name}
          </h1>
          <div class="logo">
            <FullLogo />
          </div>
        </div>
        {#if categories.length > 1}
          <wa-select
            bind:this={categorySelector}
            size="s"
            name="categoryId"
            label="Category"
            {@attach value(selectedCategoryId)}
            onchange={() => {
              selectedCategoryId = Number(categorySelector?.value);
            }}
          >
            {#each categories as category (category.id)}
              <wa-option value={category.id} label={category.name}>
                {category.name}
                <small>{formatCategoryEligibility(category)}</small>
              </wa-option>
            {/each}
          </wa-select>
        {/if}
        <div class="container" style="--num-columns: {categories.length}">
          {#each categories as category (category.id)}
            <section
              class="category"
              data-selected={category.id === selectedCategoryId
                ? "true"
                : "false"}
            >
              <Header
                {competitionId}
                categoryId={category.id}
                name={category.name}
                startTime={category.timeBegin}
                endTime={category.timeEnd}
                {scoreboard}
              ></Header>
              <ResultList
                categoryId={category.id}
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
    background-color: var(--wa-color-danger-fill-loud);
    width: 100%;
    height: 2rem;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: var(--wa-space-xs);
    color: var(--wa-color-danger-on-loud);
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
    padding: var(--wa-space-m);
  }

  .banner {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--wa-space-s);
    margin-block: var(--wa-space-s) var(--wa-space-m);
  }

  h1 {
    text-align: center;
    width: 100%;
    font-size: var(--wa-font-size-2xl);
    line-height: var(--wa-line-height-condensed);
    color: var(--wa-color-text-normal);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin: 0;
  }

  .container {
    margin-top: var(--wa-space-s);
    flex-grow: 1;
    width: 100%;
    display: grid;
    margin-inline: auto;
    max-width: calc(
      var(--num-columns) * 32rem + (var(--num-columns) - 1) * var(--wa-space-s)
    );
    grid-template-columns: repeat(var(--num-columns), 1fr);
    grid-template-rows: 1fr;
    padding-top: 0;
    gap: var(--wa-space-s);
  }

  .category {
    display: flex;
    flex-direction: column;
  }

  wa-select {
    display: none;
  }

  .logo {
    height: var(--wa-font-size-xl);
    color: var(--wa-color-text-normal);
  }

  @media screen and (max-width: 512px) {
    .banner {
      flex-direction: column-reverse;
      align-items: start;
    }

    h1 {
      font-size: var(--wa-font-size-xl);
      text-align: left;
    }

    .logo {
      height: var(--wa-font-size-l);
    }

    wa-select {
      display: block;
    }

    .container {
      grid-template-columns: 1fr;
      max-width: none;
    }

    .category[data-selected="false"] {
      display: none;
    }
  }
</style>
