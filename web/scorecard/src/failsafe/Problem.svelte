<script lang="ts">
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import SimpleTickEditor from "./SimpleTickEditor.svelte";
  import TickEditor from "./TickEditor.svelte";

  type Props = {
    problem: Problem;
    tick?: Tick;
    contenderId: number;
    usePoints: boolean;
  };

  const { problem, tick, contenderId, usePoints }: Props = $props();
</script>

<details aria-label={`Problem ${problem.number}`}>
  <summary>
    <span>
      <HoldColorIndicator
        --height="1.25rem"
        --width="1.25rem"
        primary={problem.holdColorPrimary}
        secondary={problem.holdColorSecondary}
      />
      #{problem.number}
      {#if tick?.top && tick.attemptsTop === 1}
        F
      {:else if tick?.top}
        T
      {:else if tick?.zone2}
        Z2
      {:else if tick?.zone1}
        Z1
      {/if}
    </span>
  </summary>
  {#if usePoints}
    <SimpleTickEditor {problem} {tick} {contenderId} />
  {:else}
    <TickEditor {problem} {tick} {contenderId} />
  {/if}
</details>

<style>
  summary span {
    display: inline-flex;
    align-items: center;
    gap: var(--wa-space-xs);
    vertical-align: middle;
  }
</style>
