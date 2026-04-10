<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import { type WaSelectEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import WaDropdownItem from "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import {
    EmptyState,
    Table,
    type ColumnDefinition,
  } from "@climblive/lib/components";
  import type { Raffle } from "@climblive/lib/models";
  import {
    createRaffleMutation,
    getRafflesQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { Link, navigate } from "svelte-routing";
  import DeleteRaffle from "./DeleteRaffle.svelte";

  const maxRaffles = 10;

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const rafflesQuery = $derived(getRafflesQuery(contestId));
  const createRaffle = $derived(createRaffleMutation(contestId));

  let raffles = $derived(rafflesQuery.data);

  const limitReached = $derived(
    raffles !== undefined && raffles.length >= maxRaffles,
  );

  const handleCreateRaffle = () => {
    createRaffle.mutate(undefined, {
      onSuccess: (raffle: Raffle) => navigate(`/admin/raffles/${raffle.id}`),
      onError: () => toastError("Failed to create raffle."),
    });
  };

  const columns: ColumnDefinition<Raffle>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
    },
    {
      mobile: true,
      render: renderControls,
      align: "right",
      width: "max-content",
    },
  ];
</script>

{#snippet renderName({ id }: Raffle)}
  <Link to={`/admin/raffles/${id}`}>Raffle {id}</Link>
{/snippet}

{#snippet renderControls({ id }: Raffle)}
  <DeleteRaffle raffleId={id}>
    {#snippet children({ deleteRaffle })}
      <wa-dropdown
        onwa-select={(event: WaSelectEvent) => {
          if ((event.detail.item as WaDropdownItem).value === "delete") {
            deleteRaffle();
          }
        }}
      >
        <wa-button slot="trigger" size="small" appearance="plain">
          <wa-icon name="ellipsis-vertical" label="Actions"></wa-icon>
        </wa-button>
        <wa-dropdown-item value="delete" variant="danger">
          <wa-icon slot="icon" name="trash"></wa-icon>
          Delete
        </wa-dropdown-item>
      </wa-dropdown>
    {/snippet}
  </DeleteRaffle>
{/snippet}

{#snippet createButton()}
  <wa-button
    variant="neutral"
    appearance="accent"
    onclick={handleCreateRaffle}
    disabled={limitReached}>Start new raffle</wa-button
  >
  {#if limitReached}
    <wa-callout variant="warning">
      <wa-icon slot="icon" name="circle-pause"></wa-icon>
      You have reached the maximum of {maxRaffles} raffles per contest.
    </wa-callout>
  {/if}
{/snippet}

<p class="copy">
  Raffles are used to randomly select prize winners, typically after the contest
  has ended.
</p>

<section>
  {#if raffles === undefined}
    <Loader />
  {:else if raffles.length > 0}
    {@render createButton()}
    <Table {columns} data={raffles} getId={({ id }) => id}></Table>
  {:else}
    <EmptyState
      title="No raffles yet"
      description="Start a new raffle to randomly select winners from your contenders."
    >
      {#snippet actions()}
        {@render createButton()}
      {/snippet}
    </EmptyState>
  {/if}
</section>

<style>
  section {
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: var(--wa-space-m);
  }

  .copy {
    color: var(--wa-color-text-quiet);
  }
</style>
