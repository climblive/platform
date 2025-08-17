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
    getContestQuery,
    getScoreEnginesQuery,
    startScoreEngineMutation,
    stopScoreEngineMutation,
  } from "@climblive/lib/queries";
  import { add } from "date-fns";
  import { navigate } from "svelte-routing";
  import CompClassList from "./CompClassList.svelte";
  import DuplicateContest from "./DuplicateContest.svelte";
  import ProblemList from "./ProblemList.svelte";
  import RaffleList from "./RaffleList.svelte";
  import ResultsList from "./ResultsList.svelte";
  import TicketList from "./TicketList.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let tabGroup: WaTabGroup | undefined = $state();

  const contestQuery = $derived(getContestQuery(contestId));
  const scoreEnginesQuery = $derived(getScoreEnginesQuery(contestId));
  const startScoreEngine = $derived(startScoreEngineMutation(contestId));
  const stopScoreEngine = $derived(stopScoreEngineMutation());

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
      <wa-tab slot="nav" panel="problems">Problems</wa-tab>
      <wa-tab slot="nav" panel="results">Results</wa-tab>
      <wa-tab slot="nav" panel="raffles">Raffles</wa-tab>

      <wa-tab-panel name="contest">
        <article>
          <LabeledText label="Name">
            {contest.name}
          </LabeledText>
          {#if contest.description}
            <LabeledText label="Description">
              {contest.description}
            </LabeledText>
          {/if}
          {#if contest.location}
            <LabeledText label="Location">
              {contest.location}
            </LabeledText>
          {/if}
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
        <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
        ></wa-divider>
        <CompClassList {contestId} />

        <h2>Problems</h2>
        <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
        ></wa-divider>
        <p>
          Problems are created and managed under the <a
            href="#problems"
            onclick={(e: MouseEvent) => {
              if (tabGroup) {
                tabGroup.active = "problems";
              }

              e.preventDefault();
            }}>Problems</a
          > tab.
        </p>

        <h2>Tickets</h2>
        <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
        ></wa-divider>
        <TicketList {contestId} />

        <h2>Advanced</h2>
        <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
        ></wa-divider>
        <h3>Actions</h3>
        <DuplicateContest {contestId} />
        <h3>Score Engines</h3>
        <p>
          An active score engine collects all results during a contest and
          computes scores and rankings for all participants.
        </p>
        <wa-callout variant="warning">
          <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
          <strong>Score engines are managed automatically</strong><br />
          Score engines are started automatically, and manual intervention is typically
          only required for re-scoring results long after a contest has concluded.
        </wa-callout>
        <br />

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

      <wa-tab-panel name="results">
        <h2>Results</h2>
        <ResultsList {contestId} />
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

  wa-tab-panel::part(base) {
    padding-top: var(--wa-space-s);
  }

  wa-tab-panel[name="contest"] h2 {
    margin-top: var(--wa-space-2xl);
  }
</style>
