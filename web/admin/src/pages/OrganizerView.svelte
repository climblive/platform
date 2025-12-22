<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/card/card.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import "@awesome.me/webawesome/dist/components/tag/tag.js";
  import { value } from "@climblive/lib/forms";
  import {
    getSelfQuery,
    getUsersByOrganizerQuery,
  } from "@climblive/lib/queries";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";
  import ContestList from "./ContestList.svelte";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const selectedOrganizerId =
    getContext<Writable<number | undefined>>("selectedOrganizer");

  let showAll = $state(false);
  let showAllToggle: WaSwitch | undefined = $state();

  const selfQuery = $derived(getSelfQuery());
  const usersQuery = $derived(
    getUsersByOrganizerQuery($selectedOrganizerId ?? 0, {
      enabled: !!$selectedOrganizerId,
    }),
  );

  const self = $derived(selfQuery.data);
  const users = $derived(usersQuery.data);

  const collaborators = $derived(users?.filter(({ id }) => id !== self?.id));

  let select: WaSelect | undefined = $state();

  onMount(() => {
    $selectedOrganizerId = organizerId;
  });

  const handleChange = () => {
    if (select) {
      const organizerId = Number(select.value);
      $selectedOrganizerId = organizerId;
      navigate(`/admin/organizers/${organizerId}`);
    }
  };

  const toggleShowAll = () => {
    if (!showAllToggle) {
      return;
    }

    showAll = showAllToggle.checked;
  };
</script>

{#if self}
  <div class="controls">
    {#if !showAll}
      <wa-select
        bind:this={select}
        size="small"
        appearance="outlined filled"
        {@attach value($selectedOrganizerId)}
        onchange={handleChange}
      >
        <wa-icon name="id-badge" slot="start"></wa-icon>
        {#each self.organizers as organizer (organizer.id)}
          <wa-option value={organizer.id}>{organizer.name}</wa-option>
        {/each}
      </wa-select>
    {/if}

    {#if self?.admin}
      <wa-switch bind:this={showAllToggle} onchange={toggleShowAll}
        >Show all</wa-switch
      >
    {/if}
  </div>
{/if}

<wa-card appearance="filled">
  {#if collaborators && collaborators.length > 0}
    The contests under this organizer are being managed by you and the users
    {#each collaborators as collaborator, index (collaborator.id)}
      <strong>{collaborator.username}</strong>
      {#if index === collaborators.length - 2}
        &nbsp;and&nbsp;
      {:else if index < collaborators.length - 2}
        ,&nbsp;
      {/if}
    {/each}.
  {:else}
    The contests under this organizer are being managed solely by your self.
  {/if}

  <br />

  <div class="invite-controls">
    <wa-button
      appearance="outlined"
      size="small"
      onclick={() => navigate(`./organizers/${organizerId}/invites`)}
      >Invite a friend
      <wa-icon name="user-plus" slot="start"></wa-icon>
    </wa-button>

    <wa-button
      appearance="outlined"
      size="small"
      onclick={() => navigate(`./organizers/${organizerId}/invites`)}
      >View invites
      <wa-badge variant="neutral" pill>?</wa-badge>
    </wa-button>
  </div>
</wa-card>

<ContestList organizerId={showAll ? undefined : Number(organizerId)} />

<style>
  .invite-controls {
    display: flex;
    gap: var(--wa-space-xs);
    align-items: center;
    margin-top: var(--wa-space-m);
  }

  .controls {
    display: flex;
    gap: var(--wa-space-m);
    align-items: center;

    & wa-select {
      width: 100%;
    }

    & wa-switch {
      flex-shrink: 0;
    }
  }

  wa-card {
    margin-top: var(--wa-space-m);
  }
</style>
