<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
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

  const selfQuery = $derived(getSelfQuery());
  const self = $derived($selfQuery.data);

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
</script>

<h1>Welcome!</h1>

{#if self && self.organizers.length > 1}
  <p>You can manage contests as multiple organizers.</p>

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

<ContestList organizerId={Number(organizerId)} />
