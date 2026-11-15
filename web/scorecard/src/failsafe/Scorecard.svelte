<script lang="ts">
  import {
    getCompetitorQuery,
    getCompetitionQuery,
  } from "@climblive/lib/queries";
  import EditProfile from "./EditProfile.svelte";
  import ProblemList from "./ProblemList.svelte";

  type Props = {
    competitionId: string;
    competitorId: number;
  };

  const { competitionId, competitorId }: Props = $props();

  const competitorQuery = $derived(getCompetitorQuery(competitorId));
  const competitionQuery = $derived(getCompetitionQuery(competitionId));

  const competitor = $derived(competitorQuery.data);
  const competition = $derived(competitionQuery.data);
</script>

<h2>Profile</h2>
<EditProfile {competitionId} {competitorId} />

{#if competitor?.entered && competition}
  <h2>Scorecard</h2>
  <ProblemList
    {competitionId}
    {competitorId}
    enableAttempts={!competition.usePoints ||
      competition.maxAttempts > 0 ||
      competition.pointDeduction > 0}
  />
{/if}
