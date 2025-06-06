<script lang="ts">
  import { value } from "@climblive/lib/forms";
  import { getSelfQuery } from "@climblive/lib/queries";
  import type { SlSelect } from "@shoelace-style/shoelace";
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import { type Writable } from "svelte/store";

  const selfQuery = $derived(getSelfQuery());

  const self = $derived($selfQuery.data);

  const selectedOrganizer = getContext<Writable<number>>("selectedOrganizer");

  let select: SlSelect | undefined = $state();

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
    <sl-select
      bind:this={select}
      size="small"
      {@attach value($selectedOrganizer)}
      onsl-change={handleChange}
    >
      {#each self.organizers as organizer (organizer.id)}
        <sl-option value={organizer.id}>{organizer.name}</sl-option>
      {/each}
    </sl-select>
  {/if}
</header>

<style>
  header {
    display: flex;
    align-items: center;
    justify-content: end;
    padding: 0.5rem;
    background-color: var(--sl-color-primary-600);
    height: 3rem;
  }
</style>
