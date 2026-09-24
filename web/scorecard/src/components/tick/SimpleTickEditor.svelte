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

  let {
    problem,
    tick,
    pointValue,
    putTick,
    // eslint-disable-next-line no-useless-assignment
    open = $bindable(),
  }: Props = $props();
  let latestLocalRevision = $state(0);

  const handleTick = (feature: "zone1" | "zone2" | "top", flash: boolean) => {
    navigator.vibrate?.(50);

    latestLocalRevision =
      Math.max(latestLocalRevision, tick?.revision ?? 0) + 1;

    const nextTick: Omit<Tick, "id" | "timestamp"> = {
      revision: latestLocalRevision,
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

  const flashPossible = $derived.by(() => {
    return tick?.zone1 !== true && tick?.zone2 !== true && tick?.top !== true;
  });
</script>

<div class="horizontal">
  <TickButton
    label="Top"
    onClick={() => handleTick("top", false)}
    points={pointValue?.top}
    disabled={putTick.isPending || tick?.top === true}
    iconName="check"
  />
  <TickButton
    label="Flash"
    onClick={() => handleTick("top", true)}
    points={flashPossible ? pointValue?.top : undefined}
    bonusPoints={pointValue?.flashBonus}
    disabled={putTick.isPending || !flashPossible}
    iconName="bolt"
  />
</div>

{#if problem.zone2Enabled}
  <TickButton
    label="Zone 2"
    onClick={() => handleTick("zone2", false)}
    points={pointValue?.zone2}
    disabled={putTick.isPending || tick?.zone2 === true}
    iconName="check"
  />
{/if}

{#if problem.zone1Enabled}
  <TickButton
    label="Zone 1"
    onClick={() => handleTick("zone1", false)}
    points={pointValue?.zone1}
    disabled={putTick.isPending || tick?.zone1 === true}
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
