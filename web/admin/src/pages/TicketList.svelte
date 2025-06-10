<script lang="ts">
  import Ticket from "@/components/Ticket.svelte";
  import {
    getContendersByContestQuery,
    getContestQuery,
  } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/details/details.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/tab-group/tab-group.js";
  import "@shoelace-style/shoelace/dist/components/tab-panel/tab-panel.js";
  import "@shoelace-style/shoelace/dist/components/tab/tab.js";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  const contest = $derived($contestQuery.data);
  const contenders = $derived($contendersQuery.data);

  let printDialogOpened = $state(false);

  $effect(() => {
    if (contenders && !printDialogOpened) {
      printDialogOpened = true;

      setTimeout(() => {
        window.print();
      });
    }
  });
</script>

<main>
  {#if contest && contenders}
    {#each contenders as contender (contender.id)}
      <Ticket registrationCode={contender.registrationCode} />
    {/each}
  {/if}
</main>

<style>
  @page {
    size: a4 portrait;
    margin: 2cm;
  }
</style>
