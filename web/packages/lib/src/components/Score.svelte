<script lang="ts">
  interface Props {
    value: number;
    prefix?: string;
    hideZero?: boolean;
  }

  let { value, prefix = "", hideZero = false }: Props = $props();
</script>

<div class="counter" aria-live="polite">
  {#if !(hideZero && value === 0)}
    {prefix}{value}<span class="suffix">p</span>
  {/if}
</div>
