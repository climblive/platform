<script lang="ts">
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";

  const selectedOrganizer = getContext<Writable<number>>("selectedOrganizer");

  onMount(() => {
    const organizerId = localStorage.getItem("organizerId");
    if (!organizerId) {
      return;
    }

    setTimeout(() => {
      $selectedOrganizer = Number(organizerId);

      navigate(`/admin/organizers/${organizerId}`, {
        replace: true,
      });
    });
  });
</script>
