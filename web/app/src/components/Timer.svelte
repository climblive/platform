<script lang="ts">
  import "@shoelace-style/shoelace/dist/components/tooltip/tooltip.js";
  import { format, formatDistanceToNow, intervalToDuration } from "date-fns";
  import { onDestroy, onMount } from "svelte";

  export let endTime: Date;

  let intervalTimerId: number;

  const formatTimeRemaining = () => {
    const now = new Date();

    if (endTime.getTime() - now.getTime() <= 0) {
      return "00:00:00";
    }

    const duration = intervalToDuration({ start: now, end: endTime });

    if (duration.years || duration.months || duration.weeks || duration.days) {
      return formatDistanceToNow(endTime, {});
    } else {
      return [duration.hours, duration.minutes, duration.seconds]
        .map((val) => String(val ?? 0).padStart(2, "0"))
        .join(":");
    }
  };

  let displayValue = formatTimeRemaining();

  onMount(() => {
    intervalTimerId = setInterval(() => {
      displayValue = formatTimeRemaining();
    }, 1000);
  });

  onDestroy(() => {
    clearInterval(intervalTimerId);
  });
</script>

<sl-tooltip
  content={format(endTime, "PPP pp")}
  trigger="hover click"
  placement="top-start"
>
  {displayValue}
</sl-tooltip>

<style>
  sl-tooltip::part(base) {
    display: block;
  }
</style>
