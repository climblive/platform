<script lang="ts">
  export let value: number;
  export let prefix: string = "";
  export let hideZero: boolean = false;
</script>

<div class="counter" aria-live="polite">
  {#if !(hideZero && value === 0)}
    {prefix}{value}<span class="suffix">p</span>
  {/if}
</div>
