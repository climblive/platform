<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";

  type Props = {
    onClick: (e: MouseEvent) => void;
    label: string;
    points?: number;
    active: boolean;
  };

  const { onClick, label, points, active }: Props = $props();
</script>

<wa-checkbox
  size="s"
  checked={active}
  aria-label={label}
  onclick={(event: Event) => {
    event.preventDefault();
    onClick(event as MouseEvent);
  }}
>
  {label}
  {#if points !== undefined}
    ({points}p)
  {/if}
</wa-checkbox>
