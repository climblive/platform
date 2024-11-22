<script lang="ts">
  import { onMount } from "svelte";

  export let value: number;
  export let prefix: string = "";
  export let hideZero: boolean = false;

  let counter: HTMLDivElement | undefined;
  let prevValue: number = value;

  $: {
    const animation = counter?.animate(
      [{ "--num": prevValue }, { "--num": value }],
      {
        fill: "forwards",
        duration: 250,
      },
    );

    animation?.finished.then(() => (prevValue = value));
  }

  onMount(() => {
    counter?.style.setProperty("--num", value.toString());
  });
</script>

<div bind:this={counter} class="counter" aria-live="polite">
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
    counter-reset: num var(--num);
  }

  .counter > .suffix::before {
    content: counter(num);
  }
</style>
