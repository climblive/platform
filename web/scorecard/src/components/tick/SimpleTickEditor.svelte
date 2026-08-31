<script lang="ts">
  import type { PointValue, Problem, Tick } from "@climblive/lib/models";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { CreateMutationResult } from "@tanstack/svelte-query";
  import TickButton from "./TickButton.svelte";

  interface Props {
    problem: Problem;
    tick: Tick | undefined;
    pointValue?: PointValue;
    putTick: CreateMutationResult<
      Tick,
      Error,
      Omit<Tick, "id" | "timestamp">,
      unknown
    >;
    open: boolean;
  }

  // eslint-disable-next-line no-useless-assignment
  let { problem, pointValue, putTick, open = $bindable() }: Props = $props();

  const handleTick = (feature: "zone1" | "zone2" | "top", flash: boolean) => {
    navigator.vibrate?.(50);

    const nextTick: Omit<Tick, "id" | "timestamp"> = {
      problemId: problem.id,
      top: false,
      zone2: false,
      zone1: false,
      attemptsTop: flash ? 1 : 999,
      attemptsZone2: flash ? 1 : 999,
      attemptsZone1: flash ? 1 : 999,
    };

    switch (feature) {
      case "top":
        nextTick.top = true;
        nextTick.zone2 = true;
        nextTick.zone1 = true;
        break;
      case "zone2":
        nextTick.zone2 = true;
        nextTick.zone1 = true;
        break;
      case "zone1":
        nextTick.zone1 = true;
    }

    putTick.mutate(nextTick, {
      onError: () => {
        toastUnexpectedError("Failed to register ascent.");
      },
    });

    open = false;
  };
</script>

<div class="horizontal">
  <TickButton
    label="Top"
    onClick={() => handleTick("top", false)}
    points={pointValue?.top}
    iconName="check"
  />
  <TickButton
    label="Flash"
    onClick={() => handleTick("top", true)}
    points={pointValue?.top}
    bonusPoints={pointValue?.flashBonus}
    iconName="bolt"
  />
</div>

{#if problem.zone2Enabled}
  <TickButton
    label="Zone2"
    onClick={() => handleTick("zone2", false)}
    points={pointValue?.zone2}
    iconName="check"
  />
{/if}

{#if problem.zone1Enabled}
  <TickButton
    label="Zone1"
    onClick={() => handleTick("zone1", false)}
    points={pointValue?.zone1}
    iconName="check"
  />
{/if}

<style>
  .horizontal {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }
</style>
