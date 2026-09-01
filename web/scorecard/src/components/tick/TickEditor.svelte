<script lang="ts">
  import type { PointValue, Problem, Tick } from "@climblive/lib/models";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { CreateMutationResult } from "@tanstack/svelte-query";
  import { TickBuilder } from "../../utils/tickBuilder.svelte";
  import TickBox from "./TickBox.svelte";

  interface Props {
    problem: Problem;
    tick: Tick | undefined;
    pointValue?: PointValue;
    showPoints: boolean;
    putTick: CreateMutationResult<
      Tick,
      Error,
      Omit<Tick, "id" | "timestamp">,
      unknown
    >;
  }

  const { problem, pointValue, showPoints, putTick, ...rest }: Props = $props();

  const tickBuilder = $derived(TickBuilder.from(problem, rest.tick));
  const tick = $derived(tickBuilder.tick);

  const handleSubtractAttempt = (event: MouseEvent) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    tickBuilder.subtractAttempt();

    putTick.mutate(tickBuilder.tick, {
      onError: () => {
        toastUnexpectedError("Failed to update tick.");
      },
    });
  };

  const handleAddAttempt = (event: MouseEvent) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    tickBuilder.addAttempt();

    putTick.mutate(tickBuilder.tick, {
      onError: () => {
        toastUnexpectedError("Failed to update tick.");
      },
    });
  };

  const renderSublabel = (
    featureReached: boolean,
    attempts: number,
    showFlash: boolean,
  ) => {
    if (featureReached) {
      return undefined;
    }

    switch (true) {
      case attempts === 1 && showFlash:
        return "in 1 attempt";
      default:
        return `in ${attempts} attempts`;
    }
  };

  const handleTick = (checked: boolean, feature: "zone1" | "zone2" | "top") => {
    navigator.vibrate?.(50);

    if (checked) {
      tickBuilder.reachFeature(feature);
    } else {
      tickBuilder.unreachFeature(feature);
    }

    const nextTick = tickBuilder.tick;

    putTick.mutate(nextTick, {
      onError: () => {
        toastUnexpectedError("Failed to update tick.");
      },
    });
  };
</script>

<div class="horizontal">
  <wa-button
    size="m"
    pill
    appearance="outlined"
    onclick={(event: MouseEvent) => handleSubtractAttempt(event)}
    disabled={!tickBuilder.canSubtractAttempt()}
  >
    <wa-icon name="minus"></wa-icon>
  </wa-button>

  <div class="attempts">
    <span class="number">{tick.attemptsTop}</span>
    {tick.attemptsTop === 1 ? "attempt" : "attempts"}
  </div>

  <wa-button
    size="m"
    pill
    appearance="outlined"
    onclick={(event: MouseEvent) => handleAddAttempt(event)}
    disabled={!tickBuilder.canAddAttempt()}
  >
    <wa-icon slot="start" name="plus"></wa-icon>
    Add failed attempt
  </wa-button>
</div>

<TickBox
  label="Top"
  sublabel={renderSublabel(tick.top, (tick?.attemptsTop ?? 0) + 1, true)}
  onChange={(checked) => handleTick(checked, "top")}
  points={showPoints ? pointValue?.top : undefined}
  bonusPoints={pointValue?.flashBonus}
  checked={tick?.top}
  attempts={tick?.attemptsTop ?? 0}
/>

{#if problem.zone2Enabled}
  <TickBox
    label="Zone 2"
    sublabel={renderSublabel(tick.zone2, (tick?.attemptsZone2 ?? 0) + 1, false)}
    onChange={(checked) => handleTick(checked, "zone2")}
    points={showPoints ? pointValue?.zone2 : undefined}
    checked={tick?.zone2}
    attempts={tick?.attemptsZone2 ?? 0}
  />
{/if}

{#if problem.zone1Enabled}
  <TickBox
    label="Zone 1"
    sublabel={renderSublabel(tick.zone1, (tick?.attemptsZone1 ?? 0) + 1, false)}
    onChange={(checked) => handleTick(checked, "zone1")}
    points={showPoints ? pointValue?.zone1 : undefined}
    checked={tick?.zone1}
    attempts={tick?.attemptsZone1 ?? 0}
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
