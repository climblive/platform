<script lang="ts">
  interface Props {
    value: string;
    prefix?: string;
  }

  let { value, prefix = "" }: Props = $props();
</script>

<div class="counter" aria-live="polite">
  {prefix}{value}<span class="suffix"></span>
</div>
