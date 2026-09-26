<script lang="ts">
  import { getContenderQuery, getContestQuery } from "@climblive/lib/queries";
  import EditProfile from "./EditProfile.svelte";
  import ProblemList from "./ProblemList.svelte";

  type Props = {
    contestId: number;
    contenderId: number;
  };

  const { contestId, contenderId }: Props = $props();

  const contenderQuery = $derived(getContenderQuery(contenderId));
  const contestQuery = $derived(getContestQuery(contestId));

  const contender = $derived(contenderQuery.data);
  const contest = $derived(contestQuery.data);
</script>

<h2>Profile</h2>
<EditProfile {contestId} {contenderId} />

{#if contender?.entered && contest}
  <h2>Scorecard</h2>
  <ProblemList
    {contestId}
    {contenderId}
    enableAttempts={!contest.usePoints ||
      contest.maxAttempts > 0 ||
      contest.pointDeduction > 0}
  />
{/if}
