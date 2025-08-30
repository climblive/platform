<script lang="ts">
  import Loader from "@/components/Loader.svelte";
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

  const contestQuery = $derived(getContestQuery(contestId));
  const scoreEnginesQuery = $derived(getScoreEnginesQuery(contestId));
  const startScoreEngine = $derived(startScoreEngineMutation(contestId));
  const stopScoreEngine = $derived(stopScoreEngineMutation());

  let contest = $derived($contestQuery.data);
  let scoreEngines = $derived($scoreEnginesQuery.data);

  const startPossible = $derived(
    contest?.timeBegin &&
      contest.timeBegin.getTime() - new Date().getTime() < 60 * 60 * 1_000,
  );
</script>

<wa-callout variant="warning">
  <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
  <strong>Score engines are managed automatically</strong><br />
  Score engines are started automatically, and manual intervention is typically only
  required for re-scoring results long after a contest has concluded.
</wa-callout>
<br />

{#if scoreEngines === undefined}
  <Loader />
{:else}
  {#each scoreEngines as engineInstanceId (engineInstanceId)}
    <wa-button
      appearance="outlined"
      variant="warning"
      onclick={() => $stopScoreEngine.mutate(engineInstanceId)}
      loading={$stopScoreEngine.isPending}
      >Stop engine
      <wa-icon name="stop" slot="start"></wa-icon>
    </wa-button>
  {/each}
  {#if scoreEngines.length === 0}
    <wa-button
      appearance="outlined"
      variant="warning"
      onclick={() =>
        $startScoreEngine.mutate({
          terminatedBy: add(new Date(), { hours: 6 }),
        })}
      loading={$startScoreEngine.isPending}
      disabled={!startPossible}
      >Start engine manually
      <wa-icon name="play" slot="start"></wa-icon>
    </wa-button>
  {/if}
{/if}
