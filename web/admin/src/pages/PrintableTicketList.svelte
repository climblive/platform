<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import Ticket from "@/components/Ticket.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    getCompetitorsByCompetitionQuery,
    getCompetitionQuery,
  } from "@climblive/lib/queries";

  interface Props {
    competitionId: string;
  }

  let { competitionId }: Props = $props();

  const urlParams = new URLSearchParams(window.location.search);
  const fromId = Number(urlParams.get("from")) || undefined;
  const toId = Number(urlParams.get("to")) || undefined;

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const competitorsQuery = $derived(
    getCompetitorsByCompetitionQuery(competitionId),
  );

  const competition = $derived(competitionQuery.data);
  const allCompetitors = $derived(competitorsQuery.data);

  const competitors = $derived.by(() => {
    if (!allCompetitors) {
      return undefined;
    }

    if (fromId !== undefined && toId !== undefined) {
      return allCompetitors.filter((c) => c.id >= fromId && c.id <= toId);
    }

    return allCompetitors;
  });

  let printDialogOpened = $state(false);
  let printingFinished = $state(false);

  $effect(() => {
    if (competition && competitors && !printDialogOpened) {
      printDialogOpened = true;

      setTimeout(() => {
        window.print();
      });
    }
  });
</script>

<svelte:window onafterprint={() => (printingFinished = true)} />

<main>
  {#if !competition || !competitors}
    <Loader />
  {:else}
    <section class="print-prompt">
      {#if printingFinished}
        <h1>All done?</h1>
        <p>
          Print again, or close this tab to return to managing your competition.
        </p>
        <wa-button
          variant="neutral"
          appearance="outlined"
          onclick={() => window.close()}
          size="s"
        >
          <wa-icon slot="start" name="close"></wa-icon>
          Close
        </wa-button>
        <wa-button variant="brand" onclick={() => window.print()} size="s">
          <wa-icon slot="start" name="print"></wa-icon>
          Print again
        </wa-button>
      {:else}
        <h1>Your tickets are ready</h1>
        <p>
          If the print dialog doesn't open automatically, click Print below.
        </p>
        <wa-button variant="brand" onclick={() => window.print()} size="s">
          <wa-icon slot="start" name="print"></wa-icon>
          Print
        </wa-button>
      {/if}
    </section>
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
  .print-prompt {
    text-align: center;

    & p {
      color: var(--wa-color-text-quiet);
    }

    & wa-button + wa-button {
      margin-inline-start: var(--wa-space-xs);
    }
  }

  .tickets {
    display: none;
  }

  @media print {
    .print-prompt {
      display: none;
    }

    .tickets {
      display: block;
    }
  }

  @page {
    size: a4 portrait;
    margin: 2cm;
  }
</style>
