<script lang="ts">
  import "@awesome.me/webawesome/dist/components/relative-time/relative-time.js";
  import "@awesome.me/webawesome/dist/components/tooltip/tooltip.js";
  import { format as formatTime } from "date-fns";

  type Props = {
    time: Date;
    format?: "long" | "short" | "narrow";
  };

  const { time, format }: Props = $props();
  const id = $props.id();
</script>

<wa-tooltip for={id}>{formatTime(time, "yyyy-MM-dd HH:mm")}</wa-tooltip>
<wa-relative-time {id} date={time} format={format || "long"} sync
></wa-relative-time>
