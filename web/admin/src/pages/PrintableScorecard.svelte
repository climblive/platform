<script lang="ts">
  import PaperScorecard from "@/components/PaperScorecard.svelte";
  import PrintPrompt from "@/components/PrintPrompt.svelte";
  import {
    getCategoriesQuery,
    getCompetitionQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";

  interface Props {
    competitionId: number;
  }

  const { competitionId }: Props = $props();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const categoriesQuery = $derived(getCategoriesQuery(competitionId));
  const problemsQuery = $derived(getProblemsQuery(competitionId));

  const { competition, categories, problems } = $derived({
    competition: competitionQuery.data,
    categories: categoriesQuery.data,
    problems: problemsQuery.data,
  });

  let scorecard: PaperScorecard | undefined = $state();
</script>

<svelte:head>
  <title>{competition?.name ?? "ClimbLive"} — Paper scorecard</title>
</svelte:head>

{#if competition && categories && problems?.length}
  <PrintPrompt
    readyTitle="Your scorecard is ready"
    narrowMargins
    preparePrint={() => scorecard?.preparePrint()}
  >
    <PaperScorecard
      bind:this={scorecard}
      competitionName={competition.name}
      {categories}
      {problems}
    />
  </PrintPrompt>
{/if}
