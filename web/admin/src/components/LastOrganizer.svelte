<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";

  onMount(() => {
    const organizerId = localStorage.getItem("organizerId");
    if (!organizerId) {
      return;
    }

    navigate(`/admin/organizers/${organizerId}`, {
      replace: true,
    });
  });
</script>
