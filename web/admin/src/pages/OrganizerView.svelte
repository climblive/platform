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
  import { getSelfQuery } from "@climblive/lib/queries";
  import { getContext } from "svelte";
  import { Link, navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";
  import ContestList from "./ContestList.svelte";

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

    <Link
      to={`./organizers/${organizerId}`}
      >Organizer settings and invites
    </Link>
  </div>
{/if}

<ContestList organizerId={showAll ? undefined : Number(organizerId)} />

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-m);
    flex-direction: column;
    align-items: end;

    & wa-select {
      width: 100%;
    }

    & wa-switch {
      flex-shrink: 0;
    }
  }
</style>
