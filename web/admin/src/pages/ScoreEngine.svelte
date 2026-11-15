<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import Loader from "@/components/Loader.svelte";
  import { useEvents } from "@climblive/lib/components";
  import type { EventFilter } from "@climblive/lib/models";
  import {
    getCompetitionQuery,
    getScoreEnginesByCompetitionQuery,
    startScoreEngineMutation,
    stopScoreEngineMutation,
  } from "@climblive/lib/queries";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { add, format, isBefore, subSeconds } from "date-fns";
  import { onMount } from "svelte";

  interface Props {
    competitionId: number;
  }

  let { competitionId }: Props = $props();

  const queryClient = useQueryClient();
  const { subscribeEvents, unsubscribeEvents } = useEvents();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const scoreEnginesQuery = $derived(
    getScoreEnginesByCompetitionQuery(competitionId),
  );
  const startScoreEngine = startScoreEngineMutation(competitionId);
  const stopScoreEngine = stopScoreEngineMutation();

  let competition = $derived(competitionQuery.data);
  let scoreEngines = $derived(scoreEnginesQuery.data);

  const earliestStartTime = $derived(
    competition?.timeBegin
      ? subSeconds(competition.timeBegin.getTime(), 60 * 60)
      : undefined,
  );

  onMount(() => {
    const filters: EventFilter[] = [
      {
        competitionId,
        eventTypes: {
          SCORE_ENGINE_STARTED: {},
          SCORE_ENGINE_STOPPED: {},
        },
      },
    ];
    const queryKey = ["score-engines", { competitionId }];

    const invalidateScoreEngines = () => {
      queryClient.invalidateQueries({ queryKey });
    };

    const subscriptionId = subscribeEvents(filters, invalidateScoreEngines, {
      onOpen: invalidateScoreEngines,
    });

    return () => {
      unsubscribeEvents(subscriptionId);
    };
  });
</script>

<wa-callout variant="warning" size="s">
  <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
  <strong>Score engines are managed automatically</strong><br />
  Score engines are started automatically, and manual intervention is typically only
  required for re-scoring results long after a competition has concluded.

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
  {#each scoreEngines as scoreEngine (scoreEngine.instanceId)}
    <wa-button
      size="s"
      appearance="outlined"
      variant="danger"
      onclick={() => stopScoreEngine.mutate(scoreEngine.instanceId)}
      loading={stopScoreEngine.isPending}
      >Stop engine
      <wa-icon name="stop" slot="start"></wa-icon>
    </wa-button>
  {/each}
  {#if scoreEngines.length === 0}
    <wa-button
      size="s"
      appearance="outlined"
      variant="neutral"
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
