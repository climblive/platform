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
</script>

<p class="copy">
  Tickets contain registration codes that allow competitors to enter your
  competition. These tickets may be printed on paper and distributed on site.
</p>

<CreateTicketsDialog
  bind:this={createTicketsDialog}
  {contestId}
  {remainingCodes}
/>

<div class="actions">
  <wa-button
    size="s"
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
      <wa-button appearance="outlined" size="s"
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
