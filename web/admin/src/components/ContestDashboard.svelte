<script lang="ts">
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import { ContestStateProvider } from "@climblive/lib/components";
  import {
    getCompClassesQuery,
    getContestQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { contestStateToString } from "@climblive/lib/types";
  import { getCountryName, getFlag } from "@climblive/lib/utils";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
  };

  const { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const compClassesQuery = $derived(getCompClassesQuery(contestId));
  const problemsQuery = $derived(getProblemsQuery(contestId));

  const contest = $derived(contestQuery.data);
  const compClassCount = $derived(compClassesQuery.data?.length ?? 0);
  const problemCount = $derived(problemsQuery.data?.length ?? 0);
  const contestSpansMultipleDays = $derived.by(() => {
    if (contest?.timeBegin === undefined || contest?.timeEnd === undefined) {
      return false;
    }

    return contest.timeBegin.toDateString() !== contest.timeEnd.toDateString();
  });

  const formatDateTime = (date: Date) => format(date, "MMM d, yyyy HH:mm");
  const formatDate = (date: Date) => format(date, "MMM d, yyyy");
  const formatTime = (date: Date) => format(date, "HH:mm");
</script>

{#snippet summary()}
  {@const registeredContenders = contest?.registeredContenders ?? 0}

  {@const contenderStr = (count: number) =>
    count === 1 ? "1 contender" : `${count} contenders`}

  {@const compClassStr = (count: number) =>
    count === 1 ? "1 comp class" : `${count} comp classes`}

  {@const problemStr = (count: number) =>
    count === 1 ? "1 problem" : `${count} problems`}

  <div class="summary">
    <ul>
      <li>{contenderStr(registeredContenders)}</li>
      <li>{compClassStr(compClassCount)}</li>
      <li>{problemStr(problemCount)}</li>
    </ul>
  </div>
{/snippet}

{#if contest}
  <div class="heading">
    <h2>{contest.name}</h2>
    <wa-button
      size="small"
      appearance="plain"
      onclick={() => navigate(`/admin/contests/${contest.id}/edit`)}
    >
      <wa-icon name="pencil"></wa-icon>
    </wa-button>
  </div>

  <wa-divider></wa-divider>

  {#if contest.description}
    <p class="description">{contest.description}</p>
  {/if}
  <div class="meta">
    <span>
      <wa-icon name="location-dot"></wa-icon>
      {#if contest.location}
        {contest.location},
      {/if}
      {getCountryName(contest.country)}
      {getFlag(contest.country)}
    </span>

    {#if contest.timeBegin && !contestSpansMultipleDays}
      <span>
        <wa-icon name="calendar"></wa-icon>
        {formatDate(new Date(contest.timeBegin))}
      </span>
    {/if}

    {#if contest.timeBegin && contest.timeEnd}
      <span>
        <wa-icon name="clock"></wa-icon>
        {#if contestSpansMultipleDays}
          {formatDateTime(new Date(contest.timeBegin))} – {formatDateTime(
            new Date(contest.timeEnd),
          )}
        {:else}
          {formatTime(new Date(contest.timeBegin))} – {formatTime(
            new Date(contest.timeEnd),
          )}
        {/if}
        <ContestStateProvider
          startTime={contest.timeBegin}
          endTime={contest.timeEnd}
        >
          {#snippet children({ contestState })}
            <wa-badge pill>{contestStateToString(contestState)}</wa-badge>
          {/snippet}
        </ContestStateProvider>
      </span>
    {/if}
  </div>

  <div class="summary">{@render summary()}</div>
{/if}

<style>
  wa-divider {
    --color: var(--wa-color-brand-fill-normal);
  }

  .heading {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-m);

    & > h2 {
      margin: 0;
    }
  }

  .description {
    margin-block-start: var(--wa-space-xs);
  }

  .meta {
    display: flex;
    flex-wrap: wrap;
    gap: var(--wa-space-xs);
    margin-block-start: var(--wa-space-xs);
  }

  .summary {
    margin-block-start: var(--wa-space-m);

    & ul {
      display: flex;
      list-style: none;
      margin: 0;
      padding: 0;
      flex-wrap: wrap;
    }

    & li {
      white-space: nowrap;
    }

    & li:not(:last-of-type)::after {
      content: "●";
      margin-inline: var(--wa-space-xs);
    }
  }
</style>
