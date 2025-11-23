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
