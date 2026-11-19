<script lang="ts">
  import { localizedPointsValue } from "../utils/points";
  import Points from "./Points.svelte";

  interface Props {
    value: string;
    prefix?: string;
    compactThreshold?: number;
    animate?: boolean;
  }

  let {
    value,
    prefix = "",
    compactThreshold = Number.MAX_SAFE_INTEGER,
    animate = false,
  }: Props = $props();

  const points = $derived.by(() => {
    const matches = value.match(/^([0-9]+(\.[0-9]+)?)p$/);
    if (matches === null) {
      return undefined;
    }

    return Number(matches[1]);
  });
</script>

<span>
  {prefix}{#if points === undefined}{value}{:else if animate}<Points
      value={points}
      {compactThreshold}
    />{:else}{localizedPointsValue(points, compactThreshold)}{/if}
</span>
