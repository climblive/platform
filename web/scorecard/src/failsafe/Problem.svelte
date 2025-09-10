<script lang="ts">
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import {
    createTickMutation,
    deleteTickMutation,
  } from "@climblive/lib/queries";
  import BoltIcon from "./BoltIcon.svelte";
  import DoubleCheckIcon from "./DoubleCheckIcon.svelte";

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

  const addTick = (type: "top" | "flash") => () => {
    switch (type) {
      case "top":
        $createTick.mutate({
          problemId: problem.id,
          top: true,
          attemptsTop: 999,
          zone: true,
          attemptsZone: 999,
        });

        break;
      case "flash":
        $createTick.mutate({
          problemId: problem.id,
          top: true,
          attemptsTop: 1,
          zone: true,
          attemptsZone: 1,
        });

        break;
    }
  };

  const removeTick = () => {
    if (tick?.id) {
      $deleteTick.mutate(tick.id);
    }
  };
</script>

<div class="problem" data-tick={tickType(tick)}>
  <span>
    <HoldColorIndicator
      --height="1.25rem"
      --width="1.25rem"
      primary={problem.holdColorPrimary}
      secondary={problem.holdColorSecondary}
    />
    № {problem.number}
    {#if tick?.top === true && tick?.attemptsTop === 1}
      <BoltIcon />
    {:else if tick?.top === true}
      <DoubleCheckIcon />
    {/if}
  </span>
  {#if tick}
    <button
      onclick={removeTick}
      class="wa-danger wa-small wa-pill"
      disabled={$deleteTick.isPending}>Unsend</button
    >
  {:else}
    <button
      onclick={addTick("top")}
      class="wa-danger wa-small wa-pill"
      disabled={$createTick.isPending}>Top</button
    >
    <button
      onclick={addTick("flash")}
      class="wa-danger wa-small wa-pill"
      disabled={$createTick.isPending}>Flash</button
    >
  {/if}
</div>

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
    }

    &[data-tick="flash"] {
      border-color: var(--wa-color-yellow-50);
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
