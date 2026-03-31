<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import Ticket from "@/components/Ticket.svelte";
  import {
    getContendersByContestQuery,
    getContestQuery,
  } from "@climblive/lib/queries";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const contest = $derived(contestQuery.data);
  const contenders = $derived(contendersQuery.data);

  const ticketsPerPage = 7;

  const ticketNumber = (id: number) => id.toString().padStart(6, "0");

  const chunks = $derived.by(() => {
    if (!contenders) {
      return undefined;
    }

    const result = [];

    for (let i = 0; i < contenders.length; i += ticketsPerPage) {
      result.push(contenders.slice(i, i + ticketsPerPage));
    }

    return result;
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

<main>
  {#if !contest || !chunks}
    <Loader />
  {:else}
    {#each chunks as chunk, i (i)}
      <div class="chunk" class:break-before={i > 0}>
        <span class="page-range">
          {ticketNumber(chunk[0].id)} – {ticketNumber(
            chunk[chunk.length - 1].id,
          )}
        </span>
        {#each chunk as contender (contender.id)}
          <Ticket
            contestName={contest.name}
            registrationCode={contender.registrationCode}
            ticketNumber={ticketNumber(contender.id)}
          />
        {/each}
      </div>
    {/each}
  {/if}
</main>

<style>
  @page {
    size: a4 portrait;
    margin: 2cm;
  }

  .chunk {
    position: relative;
  }

  .break-before {
    break-before: page;
  }

  .page-range {
    position: absolute;
    top: -1.5cm;
    left: 0;
    right: 0;
    text-align: center;
    font-size: var(--wa-font-size-s);
    color: var(--wa-color-text-quiet);
    font-family: monospace;
  }
</style>
