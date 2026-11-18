<script lang="ts">
  import "@awesome.me/webawesome/dist/components/number-ticker/number-ticker.js";
  import { onMount } from "svelte";
  import { supportsNumberTicker } from "../utils/numberTicker";
  import { pointsNumberFormat, roundPoints } from "../utils/points";

  interface Props {
    value: number;
    compactThreshold: number;
  }

  const { value, compactThreshold }: Props = $props();

  const rounded = $derived(roundPoints(value));
  const formatter = $derived(pointsNumberFormat(value, compactThreshold));
  const finalText = $derived(formatter.format(rounded) + "p");
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

  let supported = $state(false);

  onMount(() => {
    if (!supportsNumberTicker()) {
      return;
    }

    const motion = matchMedia("(prefers-reduced-motion: reduce)");
    const update = () => {
      supported = !motion.matches;
    };
    update();
    motion.addEventListener?.("change", update);
    return () => motion.removeEventListener?.("change", update);
  });
</script>

{#if supported && Number.isFinite(rounded)}
  <wa-number-ticker value={rounded} {valueFormatter} duration={1_000} eager>
    {finalText}
  </wa-number-ticker>
{:else}
  {finalText}
{/if}
