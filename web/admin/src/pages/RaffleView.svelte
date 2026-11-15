<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    CompetitorName,
    EmptyState,
    Table,
    useEvents,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { EventFilter, RaffleWinner } from "@climblive/lib/models";
  import {
    drawRaffleWinnerMutation,
    getCompetitorsByCompetitionQuery,
    getCompetitionQuery,
    getRaffleQuery,
    getRaffleWinnersQuery,
  } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { AxiosError } from "axios";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    raffleId: number;
  }

  let { raffleId }: Props = $props();

  const queryClient = useQueryClient();
  const { subscribeEvents, unsubscribeEvents } = useEvents();

  const raffleQuery = $derived(getRaffleQuery(raffleId));
  const drawRaffleWinner = $derived(drawRaffleWinnerMutation(raffleId));
  const raffleWinnersQuery = $derived(getRaffleWinnersQuery(raffleId));

  const raffle = $derived(raffleQuery.data);
  const sortedRaffleWinners = $derived.by(() => {
    const winners = [...(raffleWinnersQuery.data ?? [])];
    winners.sort((a, b) => {
      return b.timestamp.getTime() - a.timestamp.getTime();
    });

    return winners;
  });

  const competitionQuery = $derived(
    raffle?.competitionId
      ? getCompetitionQuery(raffle.competitionId)
      : undefined,
  );
  const competition = $derived(competitionQuery?.data);

  const competitorsQuery = $derived(
    raffle?.competitionId
      ? getCompetitorsByCompetitionQuery(raffle.competitionId)
      : undefined,
  );

  const eligibleCount = $derived.by(() => {
    const competitors = competitorsQuery?.data;

    if (competitors === undefined) {
      return undefined;
    }

    return competitors.filter(({ entered, disqualified, scrubbedAt }) => {
      switch (true) {
        case entered === undefined:
          return false;
        case disqualified:
          return false;
        case scrubbedAt !== undefined:
          return false;
        default:
          return true;
      }
    }).length;
  });

  const winnersCount = $derived(raffleWinnersQuery.data?.length ?? 0);

  const allWinnersDrawn = $derived(
    eligibleCount !== undefined && winnersCount >= eligibleCount,
  );

  $effect(() => {
    const competitionId = raffle?.competitionId;

    if (competitionId === undefined) {
      return;
    }

    const filters: EventFilter[] = [
      {
        competitionId,
        eventTypes: {
          COMPETITOR_ENTERED: {},
          COMPETITOR_DISQUALIFIED: {},
          COMPETITOR_REQUALIFIED: {},
        },
      },
    ];

    const invalidateCompetitors = () => {
      queryClient.invalidateQueries({
        queryKey: ["competitors", { competitionId }],
      });
    };

    const subscriptionId = subscribeEvents(filters, (eventName) => {
      if (
        [
          "COMPETITOR_ENTERED",
          "COMPETITOR_DISQUALIFIED",
          "COMPETITOR_REQUALIFIED",
        ].includes(eventName)
      ) {
        invalidateCompetitors();
      }
    });

    return () => {
      unsubscribeEvents(subscriptionId);
    };
  });

  const handleDrawWinner = () => {
    drawRaffleWinner.mutate(undefined, {
      onError: (error) => {
        if (error instanceof AxiosError && error.status === 404) {
          toastUnexpectedError("All winners have been drawn.");
        } else {
          toastUnexpectedError("Failed to draw winner.");
        }
      },
    });
  };

  const columns: ColumnDefinition<RaffleWinner>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "1fr",
    },
    {
      label: "Timestamp",
      mobile: true,
      render: renderTimestamp,
      width: "max-content",
    },
  ];
</script>

{#snippet renderName({
  competitorName,
  competitorId,
  competitorScrubbedAt,
}: RaffleWinner)}
  <CompetitorName
    id={competitorId}
    name={competitorName}
    scrubbedAt={competitorScrubbedAt}
    withTooltip
  />
{/snippet}

{#snippet renderTimestamp({ timestamp }: RaffleWinner)}
  {format(timestamp, "yyyy-MM-dd HH:mm")}
{/snippet}

{#snippet drawButton()}
  {#if allWinnersDrawn}
    <wa-callout variant="neutral">
      <wa-icon slot="icon" name="circle-check"></wa-icon>
      All eligible winners have been drawn.
    </wa-callout>
  {:else}
    <wa-button variant="neutral" onclick={handleDrawWinner}
      >Draw winner at random</wa-button
    >
  {/if}
{/snippet}

{#if competition && raffle}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(
          `/admin/organizers/${competition.ownership.organizerId}/competitions`,
        )}><wa-icon name="home"></wa-icon></wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() => navigate(`/admin/competitions/${raffle.competitionId}`)}
      >{competition.name}</wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() =>
        navigate(`/admin/competitions/${raffle.competitionId}#raffles`)}
      >Raffles</wa-breadcrumb-item
    >
  </wa-breadcrumb>

  <h1>Raffle {raffle.id}</h1>
  <section>
    {#if sortedRaffleWinners === undefined}
      <Loader />
    {:else if sortedRaffleWinners.length > 0}
      {@render drawButton()}
      <Table
        {columns}
        data={sortedRaffleWinners}
        getId={({ competitorId }) => competitorId}
      ></Table>
    {:else if eligibleCount === 0}
      <EmptyState
        title="No winners yet"
        description="There are no eligible winners to draw."
      />
    {:else}
      <EmptyState
        title="No winners yet"
        description="Draw the first winner of your prize raffle."
      >
        {#snippet actions()}
          {@render drawButton()}
        {/snippet}
      </EmptyState>
    {/if}
  </section>
{/if}

<style>
  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
    justify-content: start;
  }
</style>
