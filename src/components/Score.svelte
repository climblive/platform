<script lang="ts">
  export let value: number;
  export let prefix: string = "";
  export let size: "small" | "large" = "small";
  export let hideZero: boolean = false;

  let prevValue: number = 0;

  let increments: number[] | undefined = [];

  $: {
    setTimeout(() => (increments = undefined));

    const currentValue = value;
    const diff = currentValue - prevValue;

    const nextIncrements = Array.from({ length: 11 }).map((_, i) =>
      Math.ceil(prevValue + (i * diff) / 10)
    );

    setTimeout(() => (increments = nextIncrements));

    prevValue = value;
  }
</script>

<div>
  {#if increments}
    <div class="counter" data-large={size === "large"}>
      <div class="increments">
        {#each increments as inc}
          <div class={inc === 0 && hideZero ? "hidden" : undefined}>
            {prefix}{inc}p
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>

<style>
  .counter {
    --height: 1.25rem;

    height: var(--height);
    overflow: hidden;
  }

  .counter[data-large="true"] {
    --height: 2.75rem;

    font-size: var(--sl-font-size-2x-large);
  }

  .increments {
    height: 100%;
    animation: counter 0.4s forwards;
    animation-timing-function: steps(10);
  }

  .increments > * {
    height: var(--height);
  }

  .hidden {
    visibility: hidden;
  }

  @keyframes counter {
    0% {
      margin-top: 0rem;
    }
    100% {
      margin-top: calc(-10 * var(--height));
    }
  }
</style>
