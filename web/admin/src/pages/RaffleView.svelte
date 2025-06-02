<script lang="ts">
  import {
    drawRaffleWinnerMutation,
    getRaffleQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";

  interface Props {
    raffleId: number;
  }

  let { raffleId }: Props = $props();

  const raffleQuery = getRaffleQuery(raffleId);
  const drawRaffleWinner = drawRaffleWinnerMutation(raffleId);

  let raffle = $derived($raffleQuery.data);

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
  </section>
{/if}

<style>
  section {
    display: flex;
    gap: var(--sl-spacing-x-small);
    flex-direction: column;
  }
</style>
