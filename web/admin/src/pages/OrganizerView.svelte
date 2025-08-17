<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { getContext, onMount } from "svelte";
  import type { Writable } from "svelte/store";
  import ContestList from "./ContestList.svelte";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const selectedOrganizer = getContext<Writable<number>>("selectedOrganizer");

  onMount(() => {
    $selectedOrganizer = organizerId;
  });
</script>

<ContestList organizerId={Number(organizerId)} />
