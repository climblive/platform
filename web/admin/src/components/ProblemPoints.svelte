<script lang="ts">
  import { isDefined } from "@climblive/lib/utils";

  interface Props {
    pointsZone1?: number;
    pointsZone2?: number;
    pointsTop: number;
    flashBonus?: number;
    mobile?: boolean;
  }

  let {
    pointsZone1,
    pointsZone2,
    pointsTop,
    flashBonus,
    mobile = false,
  }: Props = $props();

  const values = $derived(
    [pointsZone1, pointsZone2, pointsTop].filter(isDefined),
  );
</script>

{#if mobile}
  {@const min = Math.min(...values)}
  {@const max = Math.max(...values)}
  {[min, max + (flashBonus ?? 0)].join(" - ")} pts
{:else}
  {values.join(" / ")} pts
{/if}
