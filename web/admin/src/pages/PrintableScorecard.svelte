<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import PaperScorecard from "@/components/PaperScorecard.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { EmptyState } from "@climblive/lib/components";
  import {
    getCategoriesQuery,
    getCompetitionQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { untrack } from "svelte";

  interface Props {
    competitionId: number;
  }

  const { competitionId }: Props = $props();

  const id = untrack(() => competitionId);
  const competitionQuery = getCompetitionQuery(id);
  const categoriesQuery = getCategoriesQuery(id);
  const problemsQuery = getProblemsQuery(id);

  const { competition, categories, problems, failed } = $derived({
    competition: competitionQuery.data,
    categories: categoriesQuery.data,
    problems: problemsQuery.data,
    failed: [
      competitionQuery.isError,
      categoriesQuery.isError,
      problemsQuery.isError,
    ].some(Boolean),
  });

  let scorecard: PaperScorecard | undefined = $state();

  const handlePrint = async () => {
    await document.fonts.ready;
    await scorecard?.preparePrint();
    window.print();
  };
</script>

<svelte:head>
  <title>{competition?.name ?? "ClimbLive"} — Paper scorecard</title>
</svelte:head>

{#if failed}
  <div class="print-prompt">
    <EmptyState
      title="Unable to load scorecard"
      description="The competition, categories and problems must load before printing."
    >
      {#snippet actions()}
        <wa-button onclick={() => window.location.reload()}>Try again</wa-button
        >
      {/snippet}
    </EmptyState>
  </div>
{:else if !competition || !categories || !problems}
  <div class="print-prompt"><Loader /></div>
{:else if problems.length === 0}
  <div class="print-prompt">
    <EmptyState
      title="No problems to print"
      description="Add problems to the competition before printing a paper scorecard."
    />
  </div>
{:else}
  <section class="print-prompt">
    <h1>Paper scorecard</h1>
    <p>
      Print as many blank copies as you need. Use A4 at 100% scale and turn off
      browser headers and footers. For multiple pages, you can print
      double-sided.
    </p>
    <wa-button variant="brand" onclick={handlePrint}>
      <wa-icon slot="start" name="print"></wa-icon>
      Print
    </wa-button>
  </section>
  <div class="print-content">
    <PaperScorecard
      bind:this={scorecard}
      competitionName={competition.name}
      {categories}
      {problems}
    />
  </div>
{/if}

<style>
  .print-prompt {
    margin-block-end: var(--wa-space-xl);
    text-align: center;

    & p {
      color: var(--wa-color-text-quiet);
    }
  }

  .print-content {
    position: absolute;
    inset-inline-start: -10000px;
    visibility: hidden;
  }

  @media print {
    .print-prompt {
      display: none;
    }

    .print-content {
      position: static;
      visibility: visible;
    }

    :global(html),
    :global(body) {
      background: white;
      color: black;
      overflow: visible;
      transform: none;
      will-change: auto;
    }
  }

  @page {
    size: a4 portrait;
    margin: 10mm;
  }
</style>
