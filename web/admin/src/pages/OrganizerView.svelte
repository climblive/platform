<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { value } from "@climblive/lib/forms";
  import { getSelfQuery } from "@climblive/lib/queries";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";
  import ContestList from "./ContestList.svelte";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  let showAll = $state(false);
  let showAllToggle: WaSwitch | undefined = $state();

  const selfQuery = $derived(getSelfQuery());
  const self = $derived(selfQuery.data);

  let select: WaSelect | undefined = $state();

  const selectedOrganizer =
    getContext<Writable<number | undefined>>("selectedOrganizer");

  onMount(() => {
    $selectedOrganizer = organizerId;
  });

  const handleChange = () => {
    if (select) {
      const organizerId = Number(select.value);
      $selectedOrganizer = organizerId;
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

{#if self && self.organizers.length > 1}
  <p>You can manage contests as multiple organizers.</p>

  <div class="controls">
    {#if self?.admin}
      <wa-switch bind:this={showAllToggle} onchange={toggleShowAll}
        >Show all</wa-switch
      >
    {/if}

    {#if !showAll}
      <wa-select
        bind:this={select}
        size="small"
        appearance="outlined filled"
        {@attach value($selectedOrganizer)}
        onchange={handleChange}
      >
        <wa-icon name="id-badge" slot="start"></wa-icon>
        {#each self.organizers as organizer (organizer.id)}
          <wa-option value={organizer.id}>{organizer.name}</wa-option>
        {/each}
      </wa-select>
    {/if}
  </div>
{/if}

<ContestList organizerId={showAll ? undefined : Number(organizerId)} />

<style>
  .controls {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);

    & > wa-switch {
      align-self: flex-end;
    }
  }
</style>
