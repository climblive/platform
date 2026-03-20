<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import {
    getContestQuery,
    getScoreEnginesQuery,
    startScoreEngineMutation,
    stopScoreEngineMutation,
  } from "@climblive/lib/queries";
  import { add, format, isBefore, subSeconds } from "date-fns";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const scoreEnginesQuery = $derived(getScoreEnginesQuery(contestId));
  const startScoreEngine = startScoreEngineMutation(contestId);
  const stopScoreEngine = stopScoreEngineMutation();

  let contest = $derived(contestQuery.data);
  let scoreEngines = $derived(scoreEnginesQuery.data);

  const earliestStartTime = $derived(
    contest?.timeBegin
      ? subSeconds(contest.timeBegin.getTime(), 60 * 60)
      : undefined,
  );
</script>

<wa-callout variant="warning" size="small">
  <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
  <strong>Score engines are managed automatically</strong><br />
  Score engines are started automatically, and manual intervention is typically only
  required for re-scoring results long after a contest has concluded.

  {#if earliestStartTime && isBefore(new Date(), earliestStartTime)}
    Manual start is available starting {format(
      earliestStartTime,
      "yyyy-MM-dd HH:mm",
    )}.
  {/if}
</wa-callout>
<br />

{#if scoreEngines === undefined}
  <Loader />
{:else}
  {#each scoreEngines as engineInstanceId (engineInstanceId)}
    <div class="engine">
      <code class="engine-id">{engineInstanceId}</code>
      <wa-button
        appearance="outlined"
        variant="warning"
        onclick={() => stopScoreEngine.mutate(engineInstanceId)}
        loading={stopScoreEngine.isPending}
        >Stop engine
        <wa-icon name="stop" slot="start"></wa-icon>
      </wa-button>
    </div>
  {/each}
  {#if scoreEngines.length === 0}
    <wa-button
      appearance="outlined"
      variant="warning"
      onclick={() =>
        startScoreEngine.mutate({
          terminatedBy: add(new Date(), { hours: 6 }),
        })}
      loading={startScoreEngine.isPending}
      disabled={earliestStartTime && isBefore(new Date(), earliestStartTime)}
      >Start engine manually
      <wa-icon name="play" slot="start"></wa-icon>
    </wa-button>
  {/if}
{/if}

<style>
  .engine {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
    margin-bottom: var(--wa-space-s);
  }

  .engine-id {
    font-size: var(--wa-font-size-s);
    color: var(--wa-color-neutral-fill-loud);
  }
</style>
