<script lang="ts">
  import "@awesome.me/webawesome/dist/components/relative-time/relative-time.js";
  import "@awesome.me/webawesome/dist/components/tooltip/tooltip.js";
  import { format } from "date-fns";

  type Props = {
    time: Date;
  };

  const { time }: Props = $props();
  const id = $props.id();
</script>

<wa-tooltip for={id}>{format(time, "yyyy-MM-dd HH:mm")}</wa-tooltip>
<wa-relative-time {id} date={time} format="long" sync></wa-relative-time>
