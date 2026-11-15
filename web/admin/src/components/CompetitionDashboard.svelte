<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import { CompetitionStateProvider } from "@climblive/lib/components";
  import {
    getCategoriesQuery,
    getCompetitionQuery,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import type { CompetitionState } from "@climblive/lib/types";
  import { getCountryName, getFlag } from "@climblive/lib/utils";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";

  type Props = {
    competitionId: number;
  };

  const { competitionId }: Props = $props();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const categoriesQuery = $derived(getCategoriesQuery(competitionId));
  const problemsQuery = $derived(getProblemsQuery(competitionId));

  const competition = $derived(competitionQuery.data);
  const categoryCount = $derived(categoriesQuery.data?.length ?? 0);
  const problemCount = $derived(problemsQuery.data?.length ?? 0);
  const competitionSpansMultipleDays = $derived.by(() => {
    if (
      competition?.timeBegin === undefined ||
      competition?.timeEnd === undefined
    ) {
      return false;
    }

    return (
      competition.timeBegin.toDateString() !==
      competition.timeEnd.toDateString()
    );
  });

  const formatDateTime = (date: Date) => format(date, "PPp");
  const formatDate = (date: Date) => format(date, "PP");
  const formatTime = (date: Date) => format(date, "p");

  const competitionStateToString = (state: CompetitionState): string => {
    switch (state) {
      case "NOT_STARTED":
        return "Upcoming";
      case "RUNNING":
        return "Ongoing";
      case "GRACE_PERIOD":
      case "ENDED":
        return "Ended";
    }
  };
</script>

{#snippet summary()}
  {@const registeredCompetitors = competition?.registeredCompetitors ?? 0}

  {@const competitorStr = (count: number) => `${count} registered`}

  {@const categoryStr = (count: number) =>
    count === 1 ? "1 category" : `${count} categories`}

  {@const problemStr = (count: number) =>
    count === 1 ? "1 problem" : `${count} problems`}

  <div class="summary">
    <ul>
      <li>{competitorStr(registeredCompetitors)}</li>
      <li>{categoryStr(categoryCount)}</li>
      <li>{problemStr(problemCount)}</li>
    </ul>
  </div>
{/snippet}

{#if competition}
  <div class="heading">
    <h2>
      {competition.name}
    </h2>
    <wa-button
      size="s"
      appearance="plain"
      onclick={() => navigate(`/admin/competitions/${competition.id}/edit`)}
    >
      <wa-icon name="pencil"></wa-icon>
    </wa-button>
  </div>

  <wa-divider></wa-divider>

  {#if competition.description}
    <p class="description">{competition.description}</p>
  {/if}

  <div class="meta">
    <span class="location">
      <wa-icon name="location-dot"></wa-icon>

      {#if competition.location}
        {competition.location},
      {/if}

      {getCountryName(competition.country)}
      <span class="flag">{getFlag(competition.country)}</span>
    </span>

    {#if competition.timeBegin && !competitionSpansMultipleDays}
      <span>
        <wa-icon name="calendar"></wa-icon>

        {formatDate(new Date(competition.timeBegin))}
      </span>
    {/if}

    {#if competition.timeBegin && competition.timeEnd}
      <span>
        <wa-icon name="clock"></wa-icon>

        {#if competitionSpansMultipleDays}
          {formatDateTime(new Date(competition.timeBegin))} – {formatDateTime(
            new Date(competition.timeEnd),
          )}
        {:else}
          {formatTime(new Date(competition.timeBegin))} – {formatTime(
            new Date(competition.timeEnd),
          )}
        {/if}
      </span>
    {/if}

    {#if competition.timeBegin && competition.timeEnd}
      <CompetitionStateProvider {competitionId}>
        {#snippet children({ competitionState })}
          <wa-badge pill>{competitionStateToString(competitionState)}</wa-badge>
        {/snippet}
      </CompetitionStateProvider>
    {/if}
  </div>

  <div class="summary">{@render summary()}</div>
{/if}

<style>
  wa-divider {
    --color: var(--wa-color-brand-fill-normal);
  }

  .location {
    width: 100%;
  }

  .flag {
    font-size: var(--wa-font-size-larger);
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

  .meta {
    display: flex;
    align-items: center;
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
      content: "•";
      margin-inline: var(--wa-space-xs);
    }
  }

  wa-badge {
    font-size: var(--wa-font-size-xs);
  }
</style>
