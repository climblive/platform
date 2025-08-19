<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import { FullLogo } from "@climblive/lib/components";
  import { value } from "@climblive/lib/forms";
  import { getSelfQuery } from "@climblive/lib/queries";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import { type Writable } from "svelte/store";

  const selfQuery = $derived(getSelfQuery());

  const self = $derived($selfQuery.data);

  const selectedOrganizer = getContext<Writable<number>>("selectedOrganizer");

  let select: WaSelect | undefined = $state();

  let print = $state(false);

  const handleChange = () => {
    if (select) {
      const organizerId = Number(select.value);
      $selectedOrganizer = organizerId;
      navigate(`/admin/organizers/${organizerId}`);
    }
  };

  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get("print") !== null) {
      print = true;
    }
  });
</script>

{#if !print}
  <header>
    <div>
      <p class="logo">
        <FullLogo />
      </p>
      {#if self && self.organizers.length > 1}
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
  </header>
{/if}

<style>
  header {
    background-color: var(--wa-color-surface-lowered);
  }

  div {
    margin: 0 auto;
    max-width: 1024px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-inline-end: var(--wa-space-m);
    height: 3.5rem;
    gap: var(--wa-space-xl);
  }

  .logo {
    text-align: left;
    height: var(--wa-font-size-xl);
    color: var(--wa-color-text-normal);
    padding-left: var(--wa-space-xs);
    flex-shrink: 0;
    margin-inline-start: var(--wa-space-xs);
  }

  @media print {
    header {
      display: none;
    }
  }
</style>
