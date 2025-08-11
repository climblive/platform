<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { RaffleWinner } from "@climblive/lib/models";
  import {
    drawRaffleWinnerMutation,
    getRaffleQuery,
    getRaffleWinnersQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { format } from "date-fns";

  interface Props {
    raffleId: number;
  }

  let { raffleId }: Props = $props();

  const raffleQuery = $derived(getRaffleQuery(raffleId));
  const drawRaffleWinner = $derived(drawRaffleWinnerMutation(raffleId));
  const raffleWinnersQuery = $derived(getRaffleWinnersQuery(raffleId));

  const raffle = $derived($raffleQuery.data);
  const sortedRaffleWinners = $derived.by(() => {
    const winners = [...($raffleWinnersQuery.data ?? [])];
    winners.sort((a, b) => {
      return b.timestamp.getTime() - a.timestamp.getTime();
    });

    return winners;
  });

  const handleDrawWinner = () => {
    $drawRaffleWinner.mutate(undefined, {
      onError: () => {
        toastError("Failed to draw winner.");
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

{#if raffle}
  <section>
    <h1>Raffle {raffle.id}</h1>
    <wa-button variant="brand" onclick={handleDrawWinner}>Draw winner</wa-button
    >

    {#if sortedRaffleWinners?.length}
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
    display: flex;
    gap: var(--wa-space-xs);
    flex-direction: column;
  }
</style>
