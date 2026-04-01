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

<main>
  {#if !contest || !contenders}
    <Loader />
  {:else}
    {#each contenders as contender (contender.id)}
      <Ticket
        contestName={contest.name}
        registrationCode={contender.registrationCode}
        ticketNumber={contender.id}
      />
    {/each}
  {/if}
</main>

<style>
  @page {
    size: a4 portrait;
    margin: 2cm;
  }
</style>
