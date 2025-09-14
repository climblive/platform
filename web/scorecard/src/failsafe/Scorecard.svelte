<script lang="ts">
  import { getContenderQuery } from "@climblive/lib/queries";
  import EditProfile from "./EditProfile.svelte";
  import ProblemList from "./ProblemList.svelte";

  type Props = {
    contestId: number;
    contenderId: number;
  };

  const { contestId, contenderId }: Props = $props();

  const contenderQuery = $derived(getContenderQuery(contenderId));

  const contender = $derived($contenderQuery.data);
</script>

<h2>Profile</h2>
<EditProfile {contestId} {contenderId} />

{#if contender?.entered}
  <h2>Scorecard</h2>
  <ProblemList {contestId} {contenderId}></ProblemList>
{/if}
