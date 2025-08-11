<script lang="ts">
  import type { WaTabShowEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/scroller/scroller.js";
  import "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
  import type WaTabGroup from "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
  import "@awesome.me/webawesome/dist/components/tab-panel/tab-panel.js";
  import "@awesome.me/webawesome/dist/components/tab/tab.js";
  import { LabeledText } from "@climblive/lib/components";
  import {
    duplicateContestMutation,
    getContestQuery,
    getScoreEnginesQuery,
    startScoreEngineMutation,
    stopScoreEngineMutation,
  } from "@climblive/lib/queries";
  import { getApiUrl, toastError } from "@climblive/lib/utils";
  import { add } from "date-fns";
  import { navigate } from "svelte-routing";
  import CompClassList from "./CompClassList.svelte";
  import ContenderList from "./ContenderList.svelte";
  import ProblemList from "./ProblemList.svelte";
  import RaffleList from "./RaffleList.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let tabGroup: WaTabGroup | undefined = $state();

  const contestQuery = $derived(getContestQuery(contestId));
  const scoreEnginesQuery = $derived(getScoreEnginesQuery(contestId));
  const startScoreEngine = $derived(startScoreEngineMutation(contestId));
  const stopScoreEngine = $derived(stopScoreEngineMutation());
  const duplicateContest = $derived(duplicateContestMutation(contestId));

  let contest = $derived($contestQuery.data);
  let scoreEngines = $derived($scoreEnginesQuery.data);

  $effect(() => {
    const hash = window.location.hash.substring(1);

    if (tabGroup) {
      setTimeout(() => tabGroup?.setAttribute("active", hash));
    }
  });

  const handleTabShow = (event: WaTabShowEvent) => {
    const { name } = event.detail;
    if (name) {
      window.location.hash = name;
    }
  };

  const handleDuplicationRequest = async () => {
    if (contest) {
      $duplicateContest.mutate(undefined, {
        onSuccess: (duplicate) => {
          navigate(`/admin/contests/${duplicate.id}`);
        },
        onError: () => {
          toastError("Failed to duplicate contest.");
        },
      });
    }
  };
</script>

<main>
  {#if contest && scoreEngines}
    <wa-button
      appearance="plain"
      onclick={() =>
        navigate(`/admin/organizers/${contest.ownership.organizerId}`)}
      >Back to contests<wa-icon name="arrow-left" slot="start"
      ></wa-icon></wa-button
    >
    <h1>{contest.name}</h1>

    <wa-tab-group bind:this={tabGroup} onwa-tab-show={handleTabShow}>
      <wa-tab slot="nav" panel="contest">Contest</wa-tab>
      <wa-tab slot="nav" panel="contenders">Contenders</wa-tab>
      <wa-tab slot="nav" panel="problems">Problems</wa-tab>
      <wa-tab slot="nav" panel="raffles">Raffles</wa-tab>

      <wa-tab-panel name="contest">
        <div class="actions">
          <wa-button onclick={handleDuplicationRequest} appearance="outlined"
            >Duplicate
            <wa-icon name="copy" slot="start"></wa-icon>
          </wa-button>

          <a href={`${getApiUrl()}/contests/${contestId}/results`}>
            <wa-button appearance="outlined"
              >Download results
              <wa-icon name="download" slot="start"></wa-icon>
            </wa-button>
          </a>
        </div>

        <article>
          <LabeledText label="Description">
            {contest.description}
          </LabeledText>
          <LabeledText label="Location">
            {contest.location}
          </LabeledText>
          <LabeledText label="Finalists">
            {contest.finalists}
          </LabeledText>
          <LabeledText label="Qualifying problems">
            {contest.qualifyingProblems}
          </LabeledText>
          {#if contest.rules}
            <LabeledText label="Rules">
              <wa-scroller orientation="vertical" style="max-height: 150px;">
                {@html contest.rules}
              </wa-scroller>
            </LabeledText>
          {/if}
        </article>

        <h2>Classes</h2>
        <wa-button
          variant="brand"
          appearance="accent"
          onclick={() => navigate(`contests/${contestId}/new-comp-class`)}
          >Create</wa-button
        >
        <CompClassList {contestId} />

        <h2>Score Engines</h2>
        <wa-callout variant="warning">
          <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
          <strong>Score engines are managed automatically</strong><br />
          During an active contest, score engines are initiated and maintained without
          manual intervention. Manual initiation of a score engine is typically only
          necessary for re-scoring results after a contest has ended.
        </wa-callout>
        <br />

        {#each scoreEngines as engineInstanceId (engineInstanceId)}
          <wa-button
            variant="danger"
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
            disabled={scoreEngines.length > 0}
            >Start engine manually
            <wa-icon name="play" slot="start"></wa-icon>
          </wa-button>
        {/if}
      </wa-tab-panel>

      <wa-tab-panel name="problems">
        <h2>Problems</h2>
        <wa-button
          variant="brand"
          appearance="accent"
          onclick={() => navigate(`contests/${contestId}/new-problem`)}
          >Create</wa-button
        >
        <ProblemList {contestId} />
      </wa-tab-panel>

      <wa-tab-panel name="contenders">
        <h2>Contenders</h2>
        <ContenderList {contestId} />
      </wa-tab-panel>

      <wa-tab-panel name="raffles">
        <h2>Raffles</h2>
        <RaffleList {contestId} />
      </wa-tab-panel>
    </wa-tab-group>
  {/if}
</main>

<style>
  article {
    padding-block: var(--wa-space-m);
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
  }
</style>
