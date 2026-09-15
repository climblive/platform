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

  const tickType = $derived.by(() => {
    switch (true) {
      case tick?.top && tick.attemptsTop === 1:
        return "flash";
      case tick?.top:
        return "top";
      case tick?.zone2:
        return "zone2";
      case tick?.zone1:
        return "zone1";
      default:
        return undefined;
    }
  });
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
    gap: var(--wa-space-s);
    border: var(--wa-border-width-l) var(--wa-border-style)
      var(--wa-color-surface-border);
    padding: var(--wa-space-s);
    border-radius: var(--wa-border-radius-m);

    &[data-tick="zone1"],
    &[data-tick="zone2"],
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
