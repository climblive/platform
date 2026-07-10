<script lang="ts">
  import { NO_SCORE } from "../models/score";

  interface Props {
    value: string;
    prefix?: string;
    hideZero?: boolean;
  }

  let { value, prefix = "", hideZero = false }: Props = $props();
</script>

<div class="counter" aria-live="polite">
  {#if !(hideZero && value === NO_SCORE)}
    {prefix}{value}<span class="suffix"></span>
  {/if}
</div>
