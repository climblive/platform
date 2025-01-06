<script lang="ts">
  import Ticket from "@/components/Ticket.svelte";
  import "@/print.css";

  const fakeTickets = [...Array(100).keys()];
</script>

{#each fakeTickets as fakeTicket}
  <Ticket></Ticket>
{/each}
