<script lang="ts">
  import "@awesome.me/webawesome/dist/components/number-ticker/number-ticker.js";
  import { pointsNumberFormat, roundPoints } from "../utils/points";

  interface Props {
    value: number;
    compactThreshold: number;
  }

  const { value, compactThreshold }: Props = $props();

  const rounded = $derived(roundPoints(value));
  const formatter = $derived(pointsNumberFormat(value, compactThreshold));
  const valueFormatter = $derived.by(() => {
    const fractionDigits =
      formatter.formatToParts(rounded).find((part) => part.type === "fraction")
        ?.value.length ?? 0;
    const animatedFormatter = new Intl.NumberFormat(undefined, {
      ...formatter.resolvedOptions(),
      minimumFractionDigits: fractionDigits,
      maximumFractionDigits: fractionDigits,
    });

    return (value: number) => animatedFormatter.format(value) + "p";
  });
</script>

<wa-number-ticker value={rounded} {valueFormatter} duration={1_000} eager
></wa-number-ticker>
