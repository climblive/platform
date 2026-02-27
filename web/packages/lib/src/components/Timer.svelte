<script lang="ts">
  import { formatDistanceToNow, intervalToDuration } from "date-fns";
  import { onMount } from "svelte";
  import { SyncedTime, uuidv4 } from "../utils";

  interface Props {
    endTime: Date;
    label: string;
    align?: "left" | "right";
  }

  let { endTime, label, align = "left" }: Props = $props();

  const time = new SyncedTime(1_000);

  onMount(() => {
    time.start();

    return () => time.stop();
  });

  const labelId = uuidv4();

  const formatTimeRemaining = () => {
    const now = new Date(time.current);

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

  const displayValue = $derived(formatTimeRemaining());
</script>

<div class="timer" data-alignment={align}>
  <span role="timer" aria-live="off" aria-labelledby={labelId}
    >{displayValue}</span
  >
  <label for="" id={labelId}>{label}</label>
</div>

<style>
  .timer {
    text-align: left;

    & > * {
      display: block;
      white-space: nowrap;
    }

    & span[role="timer"] {
      font-weight: var(--wa-font-weight-bold);
    }

    & label {
      font-weight: var(--wa-font-weight-normal);
      font-size: 0.75em;
    }
  }

  .timer[data-alignment="right"] {
    text-align: right;
  }
</style>
