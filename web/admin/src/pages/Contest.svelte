<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import type { WaTabShowEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/details/details.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
  import type WaTabGroup from "@awesome.me/webawesome/dist/components/tab-group/tab-group.js";
  import "@awesome.me/webawesome/dist/components/tab-panel/tab-panel.js";
  import "@awesome.me/webawesome/dist/components/tab/tab.js";
  import { LabeledText } from "@climblive/lib/components";
  import { getContestQuery } from "@climblive/lib/queries";
  import { navigate } from "svelte-routing";
  import ArchiveContest from "./ArchiveContest.svelte";
  import CompClassList from "./CompClassList.svelte";
  import DuplicateContest from "./DuplicateContest.svelte";
  import ProblemList from "./ProblemList.svelte";
  import RaffleList from "./RaffleList.svelte";
  import RestoreContest from "./RestoreContest.svelte";
  import ResultsList from "./ResultsList.svelte";
  import ScoreEngine from "./ScoreEngine.svelte";
  import TicketList from "./TicketList.svelte";
  import TransferContest from "./TransferContest.svelte";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let tabGroup: WaTabGroup | undefined = $state();
  let problemsHeading: HTMLHeadingElement | undefined = $state();
  let compClassesHeading: HTMLHeadingElement | undefined = $state();

  const contestQuery = $derived(getContestQuery(contestId));

  const contest = $derived(contestQuery.data);

  $effect(() => {
    const hash = window.location.hash.substring(1);

    if (tabGroup) {
      if (["results", "raffles"].includes(hash)) {
        setTimeout(() => tabGroup?.setAttribute("active", hash));
      }

      let scrollElement: HTMLElement | undefined;

      switch (hash) {
        case "problems":
          scrollElement = problemsHeading;
          break;
        case "comp-classes":
          scrollElement = compClassesHeading;
          break;
      }

      if (scrollElement === undefined) {
        return;
      }

      const cb = () =>
        scrollElement.scrollIntoView({
          behavior: "instant",
          block: "start",
          inline: "nearest",
        });

      if (window.requestIdleCallback !== undefined) {
        requestIdleCallback(cb);
      } else {
        setTimeout(cb, 250);
      }
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
  {#if contest === undefined}
    <Loader />
  {:else}
    <wa-breadcrumb>
      <wa-breadcrumb-item
        onclick={() =>
          navigate(`/admin/organizers/${contest.ownership.organizerId}`)}
        ><wa-icon name="home"></wa-icon></wa-breadcrumb-item
      >
      <wa-breadcrumb-item>{contest.name}</wa-breadcrumb-item>
    </wa-breadcrumb>

    {#if contest.archived === true}
      <RestoreContest {contestId} />
    {/if}

    {#if contest.archived === false}
      <wa-tab-group bind:this={tabGroup} onwa-tab-show={handleTabShow}>
        <wa-tab slot="nav" panel="contest">Contest</wa-tab>
        <wa-tab slot="nav" panel="results">Results</wa-tab>
        <wa-tab slot="nav" panel="raffles">Raffles</wa-tab>

        <wa-tab-panel name="contest">
          <article>
            <wa-button
              size="small"
              appearance="outlined"
              onclick={() => navigate(`/admin/contests/${contest.id}/edit`)}
            >
              Edit
              <wa-icon slot="start" name="pencil"></wa-icon>
            </wa-button>
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
              <wa-details summary="Rules">
                {@html contest.rules}
              </wa-details>
            {/if}
          </article>

          <h2 bind:this={compClassesHeading}>Classes</h2>
          <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
          ></wa-divider>
          <CompClassList {contestId} />

          <h2>Tickets</h2>
          <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
          ></wa-divider>
          <TicketList {contestId} />

          <h2 bind:this={problemsHeading}>Problems</h2>
          <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
          ></wa-divider>
          <ProblemList
            {contestId}
            tableLimit={window.location.hash.substring(1) === "problems"
              ? undefined
              : 8}
          />

          <h2>Advanced</h2>
          <wa-divider style="--color: var(--wa-color-brand-fill-normal);"
          ></wa-divider>
          <h3>Actions</h3>
          <div class="actions">
            <DuplicateContest {contestId} />
            <TransferContest
              {contestId}
              organizerId={contest.ownership.organizerId}
            />
            <ArchiveContest
              {contestId}
              organizerId={contest.ownership.organizerId}
            />
          </div>
          <h3>Score Engines</h3>
          <p>
            An active score engine collects all results during a contest and
            computes scores and rankings for all participants.
          </p>

          <ScoreEngine {contestId} />
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
  {/if}
</main>

<style>
  wa-breadcrumb {
    margin-block-end: var(--wa-space-m);
    display: block;
  }

  article {
    padding-block-start: var(--wa-space-m);
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  wa-tab-panel::part(base) {
    padding-top: var(--wa-space-s);
    padding-inline: var(--wa-space-2xs);
  }

  wa-tab-panel[name="contest"] h2 {
    margin-top: var(--wa-space-2xl);
  }

  wa-details {
    margin-top: var(--wa-space-m);
  }

  .actions {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
  }
</style>
