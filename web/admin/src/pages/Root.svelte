<script lang="ts">
  import { getSelfQuery } from "@climblive/lib/queries";
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";

  const selfQuery = $derived(getSelfQuery());

  const self = $derived(selfQuery.data);

  const selectedOrganizer =
    getContext<Writable<number | undefined>>("selectedOrganizer");

  $effect(() => {
    if (self === undefined) {
      return;
    }

    if (!self.organizers.some(({ id }) => id === $selectedOrganizer)) {
      $selectedOrganizer = undefined;
    }

    if ($selectedOrganizer === undefined && self.organizers.length > 0) {
      $selectedOrganizer = self.organizers[0].id;
    }

    if ($selectedOrganizer !== undefined) {
      setTimeout(() => {
        navigate(`/admin/organizers/${$selectedOrganizer}`, {
          replace: true,
        });
      });
    }
  });
</script>
