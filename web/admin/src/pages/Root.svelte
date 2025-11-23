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
    if ($selectedOrganizer !== undefined) {
      setTimeout(() => {
        navigate(`/admin/organizers/${$selectedOrganizer}`, {
          replace: true,
        });
      });
    }
  });

  $effect(() => {
    if (
      $selectedOrganizer === undefined &&
      self !== undefined &&
      self.organizers.length > 0
    ) {
      $selectedOrganizer = self.organizers[0].id;

      return;
    }
  });
</script>
