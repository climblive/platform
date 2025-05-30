<script lang="ts">
  import { getRaffleQuery } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/button/button.js";

  interface Props {
    raffleId: number;
  }

  let { raffleId }: Props = $props();

  const raffleQuery = getRaffleQuery(raffleId);

  let raffle = $derived($raffleQuery.data);
</script>

{#if raffle}
  <section>
    <h1>Raffle {raffle.id}</h1>
    <sl-button variant="primary" onclick={() => alert("Not implemented.")}
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
