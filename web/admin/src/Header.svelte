<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import { value } from "@climblive/lib/forms";
  import { getSelfQuery } from "@climblive/lib/queries";
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import { type Writable } from "svelte/store";

  const selfQuery = $derived(getSelfQuery());

  const self = $derived($selfQuery.data);

  const selectedOrganizer = getContext<Writable<number>>("selectedOrganizer");

  let select: WaSelect | undefined = $state();

  const handleChange = () => {
    if (select) {
      const organizerId = Number(select.value);
      $selectedOrganizer = organizerId;
      navigate(`/admin/organizers/${organizerId}`);
    }
  };
</script>

<header>
  {#if self && self.organizers.length > 1}
    <wa-select
      bind:this={select}
      size="small"
      appearance="filled"
      {@attach value($selectedOrganizer)}
      onchange={handleChange}
    >
      <wa-icon name="id-badge" slot="start"></wa-icon>
      {#each self.organizers as organizer (organizer.id)}
        <wa-option value={organizer.id}>{organizer.name}</wa-option>
      {/each}
    </wa-select>
  {/if}
</header>

<style>
  header {
    display: flex;
    align-items: center;
    justify-content: end;
    padding-inline: var(--wa-space-s);
    background-color: var(--wa-color-brand-fill-normal);
    height: 3.25rem;
  }

  @media print {
    header {
      display: none;
    }
  }
</style>
