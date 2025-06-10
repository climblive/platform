<script lang="ts">
  import { LabeledText } from "@climblive/lib/components";
  import {
    duplicateContestMutation,
    getContestQuery,
    getScoreEnginesQuery,
    startScoreEngineMutation,
    stopScoreEngineMutation,
  } from "@climblive/lib/queries";
  import { getApiUrl, toastError } from "@climblive/lib/utils";
  import type { SlDetails, SlTabShowEvent } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/details/details.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/tab-group/tab-group.js";
  import type SlTabGroup from "@shoelace-style/shoelace/dist/components/tab-group/tab-group.js";
  import "@shoelace-style/shoelace/dist/components/tab-panel/tab-panel.js";
  import "@shoelace-style/shoelace/dist/components/tab/tab.js";
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

  let tabGroup: SlTabGroup | undefined = $state();
  let details: SlDetails | undefined = $state();

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
      setTimeout(() => tabGroup?.show(hash));
    }
  });

  const handleTabShow = (event: SlTabShowEvent) => {
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
    <sl-button
      variant="text"
      onclick={() =>
        navigate(`/admin/organizers/${contest.ownership.organizerId}`)}
      >Back to contests<sl-icon name="arrow-left" slot="prefix"
      ></sl-icon></sl-button
    >
    <h1>{contest.name}</h1>

    <sl-tab-group bind:this={tabGroup} onsl-tab-show={handleTabShow}>
      <sl-tab slot="nav" panel="contest">Contest</sl-tab>
      <sl-tab slot="nav" panel="contenders">Contenders</sl-tab>
      <sl-tab slot="nav" panel="problems">Problems</sl-tab>
      <sl-tab slot="nav" panel="raffles">Raffles</sl-tab>

      <sl-tab-panel name="contest">
        <sl-button onclick={handleDuplicationRequest}
          >Duplicate
          <sl-icon name="copy" slot="prefix"></sl-icon>
        </sl-button>

        <a href={`${getApiUrl()}/contests/${contestId}/results`}>
          <sl-button
            >Download results
            <sl-icon name="download" slot="prefix"></sl-icon>
          </sl-button>
        </a>

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
            <sl-details
              onsl-after-show={() =>
                details?.scrollIntoView({
                  behavior: "smooth",
                  block: "start",
                  inline: "nearest",
                })}
              bind:this={details}
              summary="Rules"
            >
              {@html contest.rules}
            </sl-details>
          {/if}
        </article>

        <h2>Score Engines</h2>
        {#each scoreEngines as engineInstanceId (engineInstanceId)}
          <div>
            <h3>{engineInstanceId}</h3>
            <sl-button
              variant="danger"
              onclick={() => $stopScoreEngine.mutate(engineInstanceId)}
              loading={$stopScoreEngine.isPending}
              >Stop
              <sl-icon name="stop" slot="prefix"></sl-icon>
            </sl-button>
          </div>
        {/each}
        {#if scoreEngines.length === 0}
          <sl-button
            onclick={() =>
              $startScoreEngine.mutate({
                terminatedBy: add(new Date(), { hours: 6 }),
              })}
            loading={$startScoreEngine.isPending}
            disabled={scoreEngines.length > 0}>Start engine</sl-button
          >
        {/if}

        <h2>Classes</h2>
        <sl-button
          variant="primary"
          onclick={() => navigate(`contests/${contestId}/new-comp-class`)}
          >Create</sl-button
        >
        <CompClassList {contestId} />
      </sl-tab-panel>

      <sl-tab-panel name="problems">
        <h2>Problems</h2>
        <sl-button
          variant="primary"
          onclick={() => navigate(`contests/${contestId}/new-problem`)}
          >Create</sl-button
        >
        <ProblemList {contestId} />
      </sl-tab-panel>

      <sl-tab-panel name="contenders">
        <h2>Contenders</h2>
        <ContenderList {contestId} />
      </sl-tab-panel>

      <sl-tab-panel name="raffles">
        <h2>Raffles</h2>
        <RaffleList {contestId} />
      </sl-tab-panel>
    </sl-tab-group>
  {/if}
</main>

<style>
  article {
    padding-block: var(--sl-spacing-medium);
    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-small);
  }
</style>
