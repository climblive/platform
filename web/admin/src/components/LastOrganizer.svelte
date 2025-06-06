<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";

  onMount(() => {
    const organizerId = localStorage.getItem("organizerId");
    if (!organizerId) {
      return;
    }

    setTimeout(() => {
      navigate(`/admin/organizers/${organizerId}`, {
        replace: true,
      });
    });
  });
</script>
