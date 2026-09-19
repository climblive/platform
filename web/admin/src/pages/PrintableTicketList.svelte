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

  $effect(() => {
    if (contest && contenders && !printDialogOpened) {
      printDialogOpened = true;

      setTimeout(() => {
        window.print();
      });
    }
  });
</script>

<svelte:window onafterprint={() => window.close()} />

<main>
  {#if !contest || !contenders}
    <Loader />
  {:else}
    <section class="print-prompt">
      <h1>Your tickets are ready</h1>
      <p>
        Print the tickets for {contest.name} and hand them out to your competitors.
        Each ticket includes a unique code to join the competition.
      </p>
      <p>
        If the print dialog doesn't open automatically, click Print below. This
        tab will close when you finish with the print dialog.
      </p>
      <wa-button variant="brand" onclick={() => window.print()}>
        <wa-icon slot="start" name="print"></wa-icon>
        Print
      </wa-button>
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
    max-width: 40rem;
    margin-inline: auto;
    padding: var(--wa-space-3xl) var(--wa-space-l);
    text-align: center;

    & h1 {
      font-size: var(--wa-font-size-2xl);
    }

    & p {
      color: var(--wa-color-text-quiet);
      margin-block: var(--wa-space-m);
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
