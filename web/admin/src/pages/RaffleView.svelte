<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { Table, TableCell, TableRow } from "@climblive/lib/components";
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
</script>

{#if raffle}
  <section>
    <h1>Raffle {raffle.id}</h1>
    <wa-button variant="primary" onclick={handleDrawWinner}
      >Draw winner</wa-button
    >

    {#if sortedRaffleWinners?.length}
      <Table columns={["Name", "Timestamp"]}>
        {#each sortedRaffleWinners as winner (winner.id)}
          <TableRow>
            <TableCell>{winner.contenderName}</TableCell>
            <TableCell>{format(winner.timestamp, "yyyy-MM-dd HH:mm")}</TableCell
            >
          </TableRow>
        {/each}
      </Table>
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
