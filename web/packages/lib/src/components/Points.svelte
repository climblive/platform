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
</script>

<span class="points">
  <span class="sizer" aria-hidden="true">
    {formatter.format(rounded)}p
  </span>
  <wa-number-ticker
    value={rounded}
    valueFormatter={(value: number) => formatter.format(value) + "p"}
    duration={500}
  ></wa-number-ticker>
</span>

<style>
  .points {
    display: inline-grid;
    font-variant-numeric: tabular-nums;

    > * {
      grid-area: 1 / 1;
      white-space: nowrap;
    }

    .sizer {
      visibility: hidden;
    }
  }
</style>
