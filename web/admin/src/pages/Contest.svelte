<script lang="ts">
  import {
    getContestQuery,
    getScoreEnginesQuery,
    startScoreEngineMutation,
    stopScoreEngineMutation,
  } from "@climblive/lib/queries";
  import { add } from "date-fns";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contestQuery = getContestQuery(contestId);
  const scoreEnginesQuery = getScoreEnginesQuery(contestId);
  const startScoreEngine = startScoreEngineMutation(contestId);
  const stopScoreEngine = stopScoreEngineMutation();

  let contest = $derived($contestQuery.data);
  let scoreEngines = $derived($scoreEnginesQuery.data);
</script>

{#if contest && scoreEngines}
  <h1>{contest.name}</h1>
  <h2>Score Engines</h2>
  {#each scoreEngines as engineInstanceId}
    <div>
      <h3>{engineInstanceId}</h3>
      <sl-button
        variant="danger"
        onclick={() => $stopScoreEngine.mutate(engineInstanceId)}
        loading={$stopScoreEngine.isPending}>Stop</sl-button
      >
    </div>
  {/each}
  <sl-button
    onclick={() =>
      $startScoreEngine.mutate({ terminatedBy: add(new Date(), { hours: 1 }) })}
    loading={$startScoreEngine.isPending}
    disabled={scoreEngines.length > 0}>Start engine</sl-button
  >
{/if}
