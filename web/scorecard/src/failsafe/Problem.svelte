<script lang="ts">
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import SimpleTickEditor from "./SimpleTickEditor.svelte";
  import TickEditor from "./TickEditor.svelte";

  type Props = {
    problem: Problem;
    tick?: Tick;
    contenderId: number;
    enablePoints: boolean;
  };

  const { problem, tick, contenderId, enablePoints }: Props = $props();

  const tickType = $derived(
    tick?.top ? (tick.attemptsTop === 1 ? "flash" : "top") : "no-top",
  );
</script>

<section aria-label={`Problem ${problem.number}`} data-tick={tickType}>
  <span class="label">
    <HoldColorIndicator
      --height="1.25rem"
      --width="1.25rem"
      primary={problem.holdColorPrimary}
      secondary={problem.holdColorSecondary}
    />
    #{problem.number}
    <div class="icon">
      {#if tick?.top && tick.attemptsTop === 1}
        F
      {:else if tick?.top}
        T
      {:else if tick?.zone2}
        Z2
      {:else if tick?.zone1}
        Z1
      {/if}
    </div>
  </span>
  <div class="editor">
    {#if enablePoints}
      <SimpleTickEditor {problem} {tick} {contenderId} />
    {:else}
      <TickEditor {problem} {tick} {contenderId} />
    {/if}
  </div>
</section>

<style>
  section {
    display: grid;
    grid-template-columns: 1fr;
    gap: var(--wa-space-xs);
    border: var(--wa-border-width-m) var(--wa-border-style)
      var(--wa-color-surface-border);
    padding: var(--wa-space-s);
    border-radius: var(--wa-border-radius-m);

    &[data-tick="top"] {
      border-color: var(--wa-color-green-50);

      & .icon {
        color: var(--wa-color-green-50);
      }
    }

    &[data-tick="flash"] {
      border-color: var(--wa-color-yellow-50);

      & .icon {
        color: var(--wa-color-yellow-50);
      }
    }
  }

  .label {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
    white-space: nowrap;
  }
</style>
