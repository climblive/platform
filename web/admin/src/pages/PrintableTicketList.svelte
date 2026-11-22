<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import PrintPrompt from "@/components/PrintPrompt.svelte";
  import Ticket from "@/components/Ticket.svelte";
  import {
    getCompetitorsByCompetitionQuery,
    getCompetitionQuery,
  } from "@climblive/lib/queries";

  interface Props {
    competitionId: number;
  }

  let { competitionId }: Props = $props();

  const urlParams = new URLSearchParams(window.location.search);
  const fromId = Number(urlParams.get("from")) || undefined;
  const toId = Number(urlParams.get("to")) || undefined;

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const competitorsQuery = $derived(
    getCompetitorsByCompetitionQuery(competitionId),
  );

  const { competition, allCompetitors } = $derived({
    competition: competitionQuery.data,
    allCompetitors: competitorsQuery.data,
  });

  const competitors = $derived.by(() => {
    if (!allCompetitors) {
      return undefined;
    }

    if (fromId !== undefined && toId !== undefined) {
      return allCompetitors.filter((c) => c.id >= fromId && c.id <= toId);
    }

    return allCompetitors;
  });
</script>

<main>
  {#if !competition || !competitors}
    <Loader />
  {:else}
    <PrintPrompt readyTitle="Your tickets are ready" />
    <div class="tickets">
      {#each competitors as competitor (competitor.id)}
        <Ticket
          competitionName={competition.name}
          registrationCode={competitor.registrationCode}
          ticketNumber={competitor.id}
        />
      {/each}
    </div>
  {/if}
</main>

<style>
  .tickets {
    display: none;
    page: tickets;
  }

  @media print {
    .tickets {
      display: block;
    }
  }

  @page tickets {
    size: a4 portrait;
    margin: 2cm;
  }
</style>
