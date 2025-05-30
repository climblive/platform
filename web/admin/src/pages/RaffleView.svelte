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

<section>
  <h1>Raffle {raffleId}</h1>
  <sl-button variant="primary" onclick={() => alert("Not implemented.")}
    >Draw winner</sl-button
  >
</section>

<style>
  section {
    display: flex;
    gap: var(--sl-spacing-x-small);
    flex-direction: column;
  }
</style>
