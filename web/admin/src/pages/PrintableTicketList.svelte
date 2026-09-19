<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import Ticket from "@/components/Ticket.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    getContendersByContestQuery,
    getContestQuery,
  } from "@climblive/lib/queries";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const urlParams = new URLSearchParams(window.location.search);
  const fromId = Number(urlParams.get("from")) || undefined;
  const toId = Number(urlParams.get("to")) || undefined;

  const contestQuery = $derived(getContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const contest = $derived(contestQuery.data);
  const allContenders = $derived(contendersQuery.data);

  const contenders = $derived.by(() => {
    if (!allContenders) {
      return undefined;
    }

    if (fromId !== undefined && toId !== undefined) {
      return allContenders.filter((c) => c.id >= fromId && c.id <= toId);
    }

    return allContenders;
  });

  let printDialogOpened = $state(false);
  let printingFinished = $state(false);

  $effect(() => {
    if (contest && contenders && !printDialogOpened) {
      printDialogOpened = true;

      setTimeout(() => {
        // window.print();
      });
    }
  });
</script>

<svelte:window onafterprint={() => (printingFinished = true)} />

<main>
  {#if !contest || !contenders}
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
      {#each contenders as contender (contender.id)}
        <Ticket
          contestName={contest.name}
          registrationCode={contender.registrationCode}
          ticketNumber={contender.id}
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
