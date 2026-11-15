<script lang="ts">
  import {
    getContenderQuery,
    getCompetitionQuery,
  } from "@climblive/lib/queries";
  import EditProfile from "./EditProfile.svelte";
  import ProblemList from "./ProblemList.svelte";

  type Props = {
    competitionId: number;
    contenderId: number;
  };

  const { competitionId, contenderId }: Props = $props();

  const contenderQuery = $derived(getContenderQuery(contenderId));
  const competitionQuery = $derived(getCompetitionQuery(competitionId));

  const contender = $derived(contenderQuery.data);
  const competition = $derived(competitionQuery.data);
</script>

<h2>Profile</h2>
<EditProfile {competitionId} {contenderId} />

{#if contender?.entered && competition}
  <h2>Scorecard</h2>
  <ProblemList
    {competitionId}
    {contenderId}
    enableAttempts={!competition.usePoints ||
      competition.maxAttempts > 0 ||
      competition.pointDeduction > 0}
  />
{/if}
