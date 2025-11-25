<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import "@awesome.me/webawesome/dist/components/breadcrumb-item/breadcrumb-item.js";
  import "@awesome.me/webawesome/dist/components/breadcrumb/breadcrumb.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { RaffleWinner } from "@climblive/lib/models";
  import {
    drawRaffleWinnerMutation,
    getContestQuery,
    getRaffleQuery,
    getRaffleWinnersQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { AxiosError } from "axios";
  import { format } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    raffleId: number;
  }

  let { raffleId }: Props = $props();

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
  let contest = $derived(contestQuery?.data);

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

{#if contest && raffle}
  <wa-breadcrumb>
    <wa-breadcrumb-item
      onclick={() =>
        navigate(`/admin/organizers/${contest.ownership.organizerId}`)}
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
    <wa-button variant="neutral" onclick={handleDrawWinner}
      >Draw winner</wa-button
    >

    {#if sortedRaffleWinners === undefined}
      <Loader />
    {:else if sortedRaffleWinners.length > 0}
      <Table
        {columns}
        data={sortedRaffleWinners}
        getId={({ contenderId }) => contenderId}
      ></Table>
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
