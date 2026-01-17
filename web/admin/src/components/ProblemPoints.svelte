<script lang="ts">
  import { isDefined } from "@climblive/lib/utils";

  interface Props {
    pointsZone1?: number;
    pointsZone2?: number;
    pointsTop: number;
  }

  let { pointsZone1, pointsZone2, pointsTop }: Props = $props();

  const values = $derived(
    [pointsZone1, pointsZone2, pointsTop].filter(isDefined),
  );
</script>

{values.join(" / ")} pts
