<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import type { Problem } from "../models";

  interface Props {
    problem: Pick<Problem, "number" | "description">;
  }

  const { problem }: Props = $props();
  const id = $props.id();
</script>

{#if problem.description}
  <wa-button size="xs" appearance="plain" pill {id}
    ><wa-icon
      name="info-circle"
      label={`Description for problem ${problem.number}`}
    ></wa-icon></wa-button
  >

  <wa-popover for={id}>
    {problem.description}
    <wa-button
      appearance="outlined"
      variant="neutral"
      size="s"
      data-popover="close">Got it!</wa-button
    >
  </wa-popover>
{/if}

<style>
  wa-popover {
    --max-width: min(300px, calc(100vw - 2 * var(--wa-space-m)));

    white-space: normal;
    overflow-wrap: anywhere;
  }

  wa-popover::part(body) {
    width: var(--max-width);
  }

  wa-popover wa-button {
    margin-block-start: var(--wa-space-m);
  }
</style>
