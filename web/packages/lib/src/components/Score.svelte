<script lang="ts">
  import { onMount } from "svelte";

  export let value: number;
  export let prefix: string = "";
  export let hideZero: boolean = false;

  let counter: HTMLDivElement | undefined;
  let prevValue: number = value;

  $: {
    counter?.style.setProperty("--num", value.toString());

    setTimeout(() => (prevValue = value), 250);
  }

  onMount(() => {
    counter?.style.setProperty("--num", value.toString());
  });
</script>

<div bind:this={counter} class="counter">
  {#if !(hideZero && value === 0 && prevValue === 0)}
    {prefix}<span class="suffix">p</span>
  {/if}
</div>

<style>
  @property --num {
    syntax: "<integer>";
    initial-value: 0;
    inherits: false;
  }

  .counter {
    counter-reset: nv var(--num);
    transition: --num 250ms steps(10);
  }

  .counter > .suffix::before {
    content: counter(nv);
  }
</style>
