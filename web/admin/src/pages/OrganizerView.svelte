<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { value } from "@climblive/lib/forms";
  import { getSelfQuery } from "@climblive/lib/queries";
  import { getContext } from "svelte";
  import { Link, navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";
  import ContestList from "./ContestList.svelte";
  import CreateOrganizer from "./CreateOrganizer.svelte";

  interface Props {
    organizerId: number;
  }

  const { organizerId }: Props = $props();

  const selectedOrganizerId =
    getContext<Writable<number | undefined>>("selectedOrganizer");

  let showAll = $state(false);
  let showAllToggle: WaSwitch | undefined = $state();

  const selfQuery = $derived(getSelfQuery());

  const self = $derived(selfQuery.data);

  let select: WaSelect | undefined = $state();

  const selectedOrganizer =
    getContext<Writable<number | undefined>>("selectedOrganizer");

  $effect(() => {
    if (self) {
      if (self.organizers.some(({ id }) => id === organizerId)) {
        $selectedOrganizer = organizerId;
      } else {
        navigate("./");
      }
    }
  });

  const handleChange = () => {
    if (select) {
      const organizerId = Number(select.value);
      $selectedOrganizerId = organizerId;
      navigate(`/admin/organizers/${organizerId}/contests`);
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
  {#if !showAll}
    <div class="organizer-select">
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

      <CreateOrganizer>
        {#snippet children({ createOrganizer })}
          <wa-button
            size="small"
            variant="neutral"
            appearance="outlined"
            onclick={createOrganizer}
          >
            <wa-icon name="plus"></wa-icon>
          </wa-button>
        {/snippet}
      </CreateOrganizer>
    </div>
  {/if}

  <div class="controls">
    <div>
      {#if self?.admin}
        <wa-switch
          size="small"
          bind:this={showAllToggle}
          onchange={toggleShowAll}>Show all</wa-switch
        >
      {/if}
    </div>

    <Link to={`./organizers/${organizerId}`}
      >Organizer settings and invites</Link
    >
  </div>
{/if}

<ContestList organizerId={showAll ? undefined : Number(organizerId)} />

<style>
  .organizer-select {
    display: flex;
    align-items: center;
    justify-content: end;
    gap: var(--wa-space-xs);
  }

  .controls {
    margin-block: var(--wa-space-m) var(--wa-space-xl);
    width: 100%;
    display: flex;
    justify-content: space-between;
    font-size: var(--wa-font-size-s);
    align-items: center;
  }

  wa-select {
    width: 100%;
  }

  wa-switch {
    flex-shrink: 0;
  }
</style>
