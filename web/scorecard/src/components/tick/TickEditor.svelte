<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import type { PointValue, Problem, Tick } from "@climblive/lib/models";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { CreateMutationResult } from "@tanstack/svelte-query";
  import { isCancel } from "axios";
  import { buildTick, TickMutator } from "../../utils/tickMutator.svelte";
  import TickBox from "./TickBox.svelte";

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
  }

  const { problem, pointValue, putTick, ...rest }: Props = $props();

  const tickMutator = $derived(TickMutator.from(problem, rest.tick));
  const tick = $derived(buildTick(problem.id, tickMutator));
  let latestLocalRevision = $state(0);

  const saveTick = async () => {
    try {
      latestLocalRevision =
        Math.max(latestLocalRevision, rest.tick?.revision ?? 0) + 1;

      await putTick.mutateAsync({ ...tick, revision: latestLocalRevision });
    } catch (error) {
      if (!isCancel(error)) {
        toastUnexpectedError("Failed to update ascent.");
      }
    }
  };

  const handleSubtractAttempt = (event: MouseEvent) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    tickMutator.subtractAttempt();

    saveTick();
  };

  const handleAddAttempt = (event: MouseEvent) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    tickMutator.addAttempt();

    saveTick();
  };

  const renderSublabel = (featureReached: boolean, attempts: number) => {
    if (featureReached) {
      return undefined;
    }

    switch (true) {
      case attempts === 1:
        return "in 1 attempt";
      case attempts > 999:
        return undefined;
      default:
        return `in ${attempts} attempts`;
    }
  };

  const handleTick = (checked: boolean, feature: "zone1" | "zone2" | "top") => {
    navigator.vibrate?.(50);

    if (checked) {
      tickMutator.reachFeature(feature);
    } else {
      tickMutator.unreachFeature(feature);
    }

    saveTick();
  };
</script>

<div class="horizontal">
  <wa-button
    size="m"
    pill
    appearance="outlined"
    onclick={(event: MouseEvent) => handleSubtractAttempt(event)}
    disabled={!tickMutator.canSubtractAttempt()}
  >
    <wa-icon name="minus" label="Subtract failed attempt"></wa-icon>
  </wa-button>

  <div class="attempts" role="status" aria-label="Attempts">
    <span class="number">{tick.attemptsTop}</span>
    {tick.attemptsTop === 1 ? "attempt" : "attempts"}
  </div>

  <wa-button
    size="m"
    pill
    appearance="outlined"
    onclick={(event: MouseEvent) => handleAddAttempt(event)}
    disabled={!tickMutator.canAddAttempt()}
  >
    <wa-icon slot="start" name="plus"></wa-icon>
    Add failed attempt
  </wa-button>
</div>

<TickBox
  label="Top"
  sublabel={renderSublabel(tick.top, (tick?.attemptsTop ?? 0) + 1)}
  onChange={(checked) => handleTick(checked, "top")}
  points={pointValue?.top}
  bonusPoints={pointValue?.flashBonus}
  checked={tick?.top}
  attempts={tick?.attemptsTop ?? 0}
  disabled={tick?.top === false && !tickMutator.canAddAttempt()}
/>

{#if problem.zone2Enabled}
  <TickBox
    label="Zone 2"
    sublabel={renderSublabel(tick.zone2, (tick?.attemptsZone2 ?? 0) + 1)}
    onChange={(checked) => handleTick(checked, "zone2")}
    points={pointValue?.zone2}
    checked={tick?.zone2}
    attempts={tick?.attemptsZone2 ?? 0}
    disabled={tick?.zone2 === false && !tickMutator.canAddAttempt()}
  />
{/if}

{#if problem.zone1Enabled}
  <TickBox
    label="Zone 1"
    sublabel={renderSublabel(tick.zone1, (tick?.attemptsZone1 ?? 0) + 1)}
    onChange={(checked) => handleTick(checked, "zone1")}
    points={pointValue?.zone1}
    checked={tick?.zone1}
    attempts={tick?.attemptsZone1 ?? 0}
    disabled={tick?.zone1 === false && !tickMutator.canAddAttempt()}
  />
{/if}

<style>
  .horizontal {
    margin-inline-start: auto;
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
    flex-wrap: wrap;
  }

  .attempts {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    text-align: center;
    line-height: var(--wa-line-height-condensed);
    width: 3rem;

    & .number {
      margin: 0;
      font-size: var(--wa-font-size-l);
      font-weight: var(--wa-font-weight-bold);
      color: var(--wa-color-text-normal);
    }
  }
</style>
