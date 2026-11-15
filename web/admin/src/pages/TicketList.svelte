<script lang="ts">
  import CreateTicketsDialog from "@/components/CreateTicketsDialog.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { getCompetitorsByCompetitionQuery } from "@climblive/lib/queries";
  import { Link } from "svelte-routing";

  const maxTickets = 500;

  interface Props {
    competitionId: number;
  }

  let { competitionId }: Props = $props();

  let createTicketsDialog: CreateTicketsDialog | undefined = $state();

  const competitorsQuery = $derived(
    getCompetitorsByCompetitionQuery(competitionId),
  );

  let competitors = $derived(competitorsQuery.data);

  let remainingCodes = $derived(
    competitors === undefined ? undefined : maxTickets - competitors.length,
  );
</script>

<p class="copy">
  Tickets contain registration codes that allow competitors to enter your
  competition. These tickets may be printed on paper and distributed on site.
</p>

<CreateTicketsDialog
  bind:this={createTicketsDialog}
  {competitionId}
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
  {#if competitors && competitors.length > 0}
    <Link to={`/admin/competitions/${competitionId}/tickets`}>
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
