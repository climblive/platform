<script lang="ts">
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import { onMount } from "svelte";
  import ContestList from "./ContestList.svelte";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  onMount(() => {
    localStorage.setItem("organizerId", String(organizerId));
  });
</script>

<ContestList organizerId={Number(organizerId)} />
