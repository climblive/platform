<script lang="ts">
  interface Props {
    value: string;
    prefix?: string;
    hideZero?: boolean;
  }

  let { value, prefix = "", hideZero = false }: Props = $props();
</script>

<div class="counter" aria-live="polite">
  {#if !(hideZero && value === "0p")}
    {prefix}{value}<span class="suffix">p</span>
  {/if}
</div>
