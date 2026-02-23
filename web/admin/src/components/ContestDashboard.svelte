<script lang="ts">
  import {
    getCompClassesQuery,
    getContestQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
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
      </span>
    {/if}
  </div>

  <div class="summary">
    {contest.registeredContenders} contender{contest.registeredContenders === 1
      ? ""
      : "s"}
    across {compClassCount} class{compClassCount === 1 ? "" : "es"}
    attempted {problemCount} problem{problemCount === 1 ? "" : "s"}.
  </div>
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
  }
</style>
