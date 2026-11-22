<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import PaperScorecard from "@/components/PaperScorecard.svelte";
  import PrintPrompt from "@/components/PrintPrompt.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
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
  <PrintPrompt
    readyTitle="Your scorecard is ready"
    preparePrint={() => scorecard?.preparePrint()}
  />
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
  }

  .print-content {
    position: absolute;
    inset-inline-start: -10000px;
    visibility: hidden;
    page: scorecard;
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

  @page scorecard {
    size: a4 portrait;
    margin: 10mm;
  }
</style>
