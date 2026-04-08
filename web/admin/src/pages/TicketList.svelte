<script lang="ts">
  import CreateTicketsDialog from "@/components/CreateTicketsDialog.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { getContendersByContestQuery } from "@climblive/lib/queries";
  import { Link } from "svelte-routing";

  const maxTickets = 500;

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let createTicketsDialog: CreateTicketsDialog | undefined = $state();

  const contendersQuery = $derived(getContendersByContestQuery(contestId));

  let contenders = $derived(contendersQuery.data);

  let remainingCodes = $derived(
    contenders === undefined ? undefined : maxTickets - contenders.length,
  );

  let registeredContenders = $derived.by(() => {
    if (!contenders) {
      return undefined;
    }

    let count = 0;

    for (const contender of contenders) {
      if (contender.entered !== undefined) {
        count += 1;
      }
    }

    return count;
  });
</script>

<p class="copy">
  Tickets contain registration codes that allow the contenders to enter your
  contest. These tickets may be printed on paper and distributed to the
  contenders on site.
  {#if contenders && contenders.length > 0}
    Out of the {contenders.length}
    tickets that you have created, {registeredContenders} have already been used.
  {/if}
</p>

<CreateTicketsDialog
  bind:this={createTicketsDialog}
  {contestId}
  {remainingCodes}
/>

<div class="actions">
  <wa-button
    size="small"
    variant="neutral"
    appearance="accent"
    onclick={() => createTicketsDialog?.open()}
    disabled={remainingCodes === undefined || remainingCodes === 0}
  >
    <wa-icon slot="start" name="plus"></wa-icon>
    Create tickets</wa-button
  >
  {#if contenders && contenders.length > 0}
    <Link to={`/admin/contests/${contestId}/tickets`}>
      <wa-button appearance="outlined" size="small"
        >View and print tickets
        <wa-icon name="list" slot="start"></wa-icon>
      </wa-button>
    </Link>
  {/if}
</div>

<p>
  {#if remainingCodes === maxTickets}
    You may create up to {maxTickets} tickets.
  {:else}
    You may create {remainingCodes} more tickets.
  {/if}
</p>

<style>
  .actions {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
  }

  .copy {
    color: var(--wa-color-text-quiet);
  }
</style>
