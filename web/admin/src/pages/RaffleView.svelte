<script lang="ts">
  import { Table, TableCell, TableRow } from "@climblive/lib/components";
  import {
    drawRaffleWinnerMutation,
    getRaffleQuery,
    getRaffleWinnersQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import { format } from "date-fns";

  interface Props {
    raffleId: number;
  }

  let { raffleId }: Props = $props();

  const raffleQuery = getRaffleQuery(raffleId);
  const drawRaffleWinner = drawRaffleWinnerMutation(raffleId);
  const raffleWinnersQuery = getRaffleWinnersQuery(raffleId);

  const raffle = $derived($raffleQuery.data);
  const sortedRaffleWinners = $derived(() => {
    const winners = [...($raffleWinnersQuery.data ?? [])];
    winners.sort((a, b) => {
      return a.timestamp.getTime() - b.timestamp.getTime();
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
    <sl-button variant="primary" onclick={handleDrawWinner}
      >Draw winner</sl-button
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
    gap: var(--sl-spacing-x-small);
    flex-direction: column;
  }
</style>
