<script lang="ts">
  import { formatDistanceToNow, intervalToDuration } from "date-fns";
  import { onDestroy, onMount } from "svelte";

  export let endTime: Date;

  let intervalTimerId: number;

  const formatTimeRemaining = () => {
    const now = new Date();

    if (endTime.getTime() - now.getTime() <= 0) {
      return "00:00:00";
    }

    now.setMilliseconds(0);

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

  const tick = () => {
    displayValue = formatTimeRemaining();

    const firefoxEarlyWakeUpCompensation = 1;

    const drift = Date.now() % 1_000;
    const next = 1_000 - drift + firefoxEarlyWakeUpCompensation;

    intervalTimerId = setTimeout(tick, next);
  };

  onMount(() => {
    tick();
  });

  onDestroy(() => {
    clearInterval(intervalTimerId);
  });
</script>

<span role="timer" aria-live="off">{displayValue}</span>
