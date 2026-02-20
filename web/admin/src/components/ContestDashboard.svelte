<script lang="ts">
  import {
    getCompClassesQuery,
    getContendersByContestQuery,
    getContestQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { getCountryName, getFlag } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
  };

  const { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const compClassesQuery = $derived(getCompClassesQuery(contestId));
  const problemsQuery = $derived(getProblemsQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const contest = $derived(contestQuery.data);
  const compClassCount = $derived(compClassesQuery.data?.length ?? 0);
  const problemCount = $derived(problemsQuery.data?.length ?? 0);
  const ticketCount = $derived(contendersQuery.data?.length ?? 0);
  const contestSpansMultipleDays = $derived.by(() => {
    if (contest?.timeBegin === undefined || contest?.timeEnd === undefined) {
      return false;
    }

    return contest.timeBegin.toDateString() !== contest.timeEnd.toDateString();
  });

  const formatDateTime = (date: Date) =>
    date.toLocaleString(undefined, {
      dateStyle: "medium",
      timeStyle: "short",
    });
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
        {contest.location}
      {/if},
      {getCountryName(contest.country)}
      {getFlag(contest.country)}
    </span>

    {#if contest.timeBegin && !contestSpansMultipleDays}
      <span>
        <wa-icon name="calendar"></wa-icon>

        {new Date(contest.timeBegin).toLocaleDateString(undefined, {
          dateStyle: "medium",
        })}
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
          {new Date(contest.timeBegin).toLocaleTimeString(undefined, {
            hour: "2-digit",
            minute: "2-digit",
          })}–{new Date(contest.timeEnd).toLocaleTimeString(undefined, {
            hour: "2-digit",
            minute: "2-digit",
          })}
        {/if}
      </span>
    {/if}
  </div>

  <div class="meta">
    <wa-icon name="users"></wa-icon>
    {contest.registeredContenders} / {ticketCount} contenders •
    {compClassCount} classes •
    {problemCount} problems
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
    font-size: var(--wa-font-size-s);
  }
</style>
