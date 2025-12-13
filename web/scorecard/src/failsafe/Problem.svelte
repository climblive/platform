<script lang="ts">
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import {
      createTickMutation,
      deleteTickMutation,
  } from "@climblive/lib/queries";

  type Props = {
    problem: Problem;
    tick?: Tick;
    contenderId: number;
  };

  const { problem, tick, contenderId }: Props = $props();

  const createTick = $derived(createTickMutation(contenderId));
  const deleteTick = $derived(deleteTickMutation());

  const tickType = (tick?: Tick) => {
    if (tick?.top && tick.attemptsTop === 1) {
      return "flash";
    }

    if (tick?.top) {
      return "top";
    }

    return "no-top";
  };

  const addTick = (type: "zone1" | "zone2" | "top" | "flash") => () => {
    switch (type) {
      case "top":
        createTick.mutate({
          problemId: problem.id,
          zone1: true,
          attemptsZone1: 999,
          zone2: true,
          attemptsZone2: 999,
          top: true,
          attemptsTop: 999,
        });

        break;
      case "flash":
        createTick.mutate({
          problemId: problem.id,
          zone1: true,
          attemptsZone1: 1,
          zone2: true,
          attemptsZone2: 1,
          top: true,
          attemptsTop: 1,
        });

        break;
    }
  };

  const removeTick = () => {
    if (tick?.id) {
      deleteTick.mutate(tick.id);
    }
  };
</script>

<section
  aria-label={`Problem ${problem.number}`}
  class="problem"
  data-tick={tickType(tick)}
>
  <span>
    <HoldColorIndicator
      --height="1.25rem"
      --width="1.25rem"
      primary={problem.holdColorPrimary}
      secondary={problem.holdColorSecondary}
    />
    № {problem.number}
    <div class="icon">
      {#if tick?.top && tick?.attemptsTop === 1}
        F
      {:else if tick?.top}
        T
      {:else if tick?.zone1}
        Z1
      {:else if tick?.zone2}
        Z2
      {/if}
    </div>
  </span>
  {#if tick}
    <button onclick={removeTick} disabled={deleteTick.isPending}>Unsend</button>
  {:else}
    {#if problem.zone1Enabled}
      <button onclick={addTick("zone1")} disabled={createTick.isPending}
        >Zone 1</button
      >
    {/if}
    {#if problem.zone2Enabled}
      <button onclick={addTick("zone2")} disabled={createTick.isPending}
        >Zone 2</button
      >
    {/if}
    <button onclick={addTick("top")} disabled={createTick.isPending}>Top</button
    >
    <button onclick={addTick("flash")} disabled={createTick.isPending}
      >Flash</button
    >
  {/if}
</section>

<style>
  .problem {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-m);
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

    & span {
      display: flex;
      align-items: center;
      gap: var(--wa-space-xs);
      white-space: nowrap;
      flex-grow: 1;
      width: max-content;

      & :global(*) {
        flex-shrink: 0;
      }
    }
  }
</style>
