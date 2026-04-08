<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { RaffleWinner } from "@climblive/lib/models";
  import {
    drawRaffleWinnerMutation,
    getContendersByContestQuery,
    getContestQuery,
    getRaffleQuery,
    getRaffleWinnersQuery,
  } from "@climblive/lib/queries";
  import { getApiUrl, toastError } from "@climblive/lib/utils";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { AxiosError } from "axios";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    raffleId: number;
  }

  let { raffleId }: Props = $props();

  const queryClient = useQueryClient();

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

  const contestQuery = $derived(
    raffle?.contestId ? getContestQuery(raffle.contestId) : undefined,
  );
  const contest = $derived(contestQuery?.data);

  const contendersQuery = $derived(
    raffle?.contestId
      ? getContendersByContestQuery(raffle.contestId)
      : undefined,
  );

  const eligibleCount = $derived.by(() => {
    const contenders = contendersQuery?.data;

    if (contenders === undefined) {
      return undefined;
    }

    return contenders.filter(
      ({ entered, disqualified }) => entered !== undefined && !disqualified,
    ).length;
  });

  const winnersCount = $derived(raffleWinnersQuery.data?.length ?? 0);

  const allWinnersDrawn = $derived(
    eligibleCount !== undefined && winnersCount >= eligibleCount,
  );

  $effect(() => {
    const contestId = raffle?.contestId;

    if (contestId === undefined) {
      return;
    }

    const eventSource = new EventSource(
      `${getApiUrl()}/contests/${contestId}/events`,
    );

    const invalidateContenders = () => {
      queryClient.invalidateQueries({
        queryKey: ["contenders", { contestId }],
      });
    };

    eventSource.addEventListener("CONTENDER_ENTERED", invalidateContenders);
    eventSource.addEventListener(
      "CONTENDER_PUBLIC_INFO_UPDATED",
      invalidateContenders,
    );
    eventSource.addEventListener(
      "CONTENDER_DISQUALIFIED",
      invalidateContenders,
    );
    eventSource.addEventListener("CONTENDER_REQUALIFIED", invalidateContenders);

    return () => {
      eventSource.close();
    };
  });

  const handleDrawWinner = () => {
    drawRaffleWinner.mutate(undefined, {
      onError: (error) => {
        if (error instanceof AxiosError && error.status === 404) {
          toastError("All winners have been drawn.");
        } else {
          toastError("Failed to draw winner.");
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

{#snippet renderName({ contenderName }: RaffleWinner)}
  {contenderName}
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
      >Draw winner {#if eligibleCount !== undefined}{winnersCount + 1} of {eligibleCount}{/if}</wa-button
    >
  {/if}
{/snippet}

{#if contest && raffle}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(`/admin/organizers/${contest.ownership.organizerId}/contests`)}
      ><wa-icon name="home"></wa-icon></wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() => navigate(`/admin/contests/${raffle.contestId}`)}
      >{contest.name}</wa-breadcrumb-item
    >
    <wa-breadcrumb-item
      onclick={() => navigate(`/admin/contests/${raffle.contestId}#raffles`)}
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
        getId={({ contenderId }) => contenderId}
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
    gap: var(--wa-space-xs);
    justify-content: start;
  }

  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }
</style>
