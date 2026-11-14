<script lang="ts">
  import { localizedPointsValue } from "../utils";

  interface Props {
    value: string;
    prefix?: string;
    compactThreshold?: number;
  }

  let {
    value,
    prefix = "",
    compactThreshold = Number.MAX_SAFE_INTEGER,
  }: Props = $props();

  const localizedValue = $derived.by(() => {
    const matches = value.match(/^([0-9]+(\.[0-9]+)?)p$/);
    if (matches === null) {
      return value;
    }

    return localizedPointsValue(Number(matches[1]), compactThreshold);
  });
</script>

<span>
  {prefix}{localizedValue}
</span>
