<script lang="ts">
  import { formatDistanceToNow, intervalToDuration } from "date-fns";
  import { onDestroy, onMount } from "svelte";
  import { uuidv4 } from "../utils";

  interface Props {
    endTime: Date;
    label: string;
  }

  let { endTime, label }: Props = $props();

  let intervalTimerId: number;
  const labelId = uuidv4();

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

  let displayValue = $state(formatTimeRemaining());

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

<div class="timer">
  <span role="timer" aria-live="off" aria-labelledby={labelId}
    >{displayValue}</span
  >
  <label for="" id={labelId}>{label}</label>
</div>

<style>
  .timer {
    text-align: right;

    & > * {
      display: block;
    }

    & span[role="timer"] {
      font-weight: var(--sl-font-weight-bold);
    }

    & label {
      font-weight: var(--sl-font-weight-normal);
      font-size: 0.75em;
    }
  }
</style>
