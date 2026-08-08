<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { PointValue, Problem, Tick } from "@climblive/lib/models";
  import { deleteTickMutation, putTickMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { AxiosError } from "axios";
  import { getContext } from "svelte";
  import type { Readable } from "svelte/store";
  import TickButton from "./TickButton.svelte";
  import { TickBuilder } from "./tick";

  interface Props {
    problem: Problem;
    tick: Tick | undefined;
    disabled: boolean | undefined;
    pointValue?: PointValue;
    showPoints: boolean;
    enableAttempts: boolean;
  }

  const {
    problem,
    tick,
    disabled = false,
    pointValue,
    showPoints,
    enableAttempts,
  }: Props = $props();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");
  const putTick = $derived(putTickMutation($session.contenderId));
  const deleteTick = $derived(deleteTickMutation());

  const tickBuilder = $derived(new TickBuilder(problem, tick));

  let open = $state(false);

  const loading = $derived(putTick.isPending || deleteTick.isPending);

  const variant = $derived.by(() => {
    switch (true) {
      case tick?.top && tick.attemptsTop === 1:
        return "flash";
      case tick?.top:
        return "top";
      case tick?.zone2 && problem.zone2Enabled:
        return "zone2";
      case tick?.zone1 && problem.zone1Enabled:
        return "zone1";
    }
  });

  const handleDelete = (event: MouseEvent) => {
    if (!tick) {
      return;
    }

    event.stopPropagation();

    open = false;

    deleteTick.mutate(tick.id, {
      onError: (error) => {
        if (error instanceof AxiosError && error.status === 404) {
          toastUnexpectedError("Ascent is already removed.");
        } else {
          toastUnexpectedError("Failed to remove tick.");
        }
      },
    });
  };

  const handleTick = (
    checked: boolean,
    feature: "zone1" | "zone2" | "top",
    flash: boolean,
  ) => {
    navigator.vibrate?.(50);

    if (checked) {
      tickBuilder.reachFeature(feature);
    } else {
      tickBuilder.unreachFeature(feature);
    }

    const nextTick = tickBuilder.tick;

    if (!enableAttempts) {
      let attempts = 999;
      if (flash && checked) {
        attempts = 1;
      }

      nextTick.attemptsTop = attempts;
      nextTick.attemptsZone2 = attempts;
      nextTick.attemptsZone1 = attempts;
    }

    putTick.mutate(nextTick, {
      onError: () => {
        toastUnexpectedError("Failed to update tick.");
      },
    });
  };

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
</script>

<div class="container">
  <button
    data-variant={variant}
    disabled={disabled || loading}
    onclick={() => (open = true)}
    aria-label={tick?.id ? "Edit" : "Tick"}
  >
    {#if loading}
      <wa-spinner></wa-spinner>
    {:else if tick?.top && tick.attemptsTop === 1}
      <wa-icon name="bolt"></wa-icon>
    {:else if tick?.top}
      T
    {:else if tick?.zone2 && problem.zone2Enabled}
      Z2
    {:else if tick?.zone1 && problem.zone1Enabled}
      Z1
    {:else if tick !== undefined && enableAttempts}
      {#if tick.attemptsTop > 99}
        <span class="small">99+</span>
      {:else}
        {tick.attemptsTop}
      {/if}
    {/if}
  </button>

  <wa-dialog
    label="Problem number {problem.number}"
    {open}
    onwa-hide={() => (open = false)}
  >
    <div class="label" slot="label">
      <HoldColorIndicator
        --height="1.25em"
        --width="1.25em"
        primary={problem.holdColorPrimary}
        secondary={problem.holdColorSecondary}
      /> Problem № {problem.number}
    </div>

    <div class="horizontal">
      <TickButton
        label="Top"
        subLabel="T"
        onChange={(checked, flash) => handleTick(checked, "top", flash)}
        points={showPoints ? pointValue?.top : undefined}
        bonusPoints={pointValue?.flashBonus}
        checked={tick?.top}
        attempts={tick?.attemptsTop ?? 0}
        flashToggle
        showAttempts={enableAttempts}
        {showPoints}
      />
    </div>

    {#if problem.zone2Enabled}
      <TickButton
        label="Zone 2"
        subLabel="Z₂"
        onChange={(checked) => handleTick(checked, "zone2", false)}
        points={showPoints ? pointValue?.zone2 : undefined}
        checked={tick?.zone2}
        indeterminate={tick?.top}
        attempts={tick?.attemptsZone2 ?? 0}
        showAttempts={enableAttempts}
        {showPoints}
      />
    {/if}

    {#if problem.zone1Enabled}
      <TickButton
        label="Zone 1"
        subLabel="Z₁"
        onChange={(checked) => handleTick(checked, "zone1", false)}
        points={showPoints ? pointValue?.zone1 : undefined}
        checked={tick?.zone1}
        indeterminate={tick?.zone2}
        attempts={tick?.attemptsZone1 ?? 0}
        showAttempts={enableAttempts}
        {showPoints}
      />
    {/if}

    {#if enableAttempts}
      <div class="horizontal">
        <wa-button
          size="s"
          appearance="outlined"
          onclick={(event: MouseEvent) => handleAddAttempt(event)}
          disabled={!tickBuilder.canAddAttempt()}
        >
          <wa-icon slot="start" name="plus"></wa-icon>
          Attempt
        </wa-button>

        <wa-button
          size="s"
          appearance="outlined"
          onclick={(event: MouseEvent) => handleSubtractAttempt(event)}
          disabled={!tickBuilder.canSubtractAttempt()}
        >
          <wa-icon slot="start" name="minus"></wa-icon>
          Attempt
        </wa-button>
      </div>
    {/if}

    <wa-button
      size="s"
      appearance="plain"
      onclick={(e: MouseEvent) => handleDelete(e)}
      variant="danger"
      disabled={open && tick === undefined}
    >
      <wa-icon slot="start" name="rotate-left"></wa-icon>
      Remove
    </wa-button>
  </wa-dialog>
</div>

<style>
  .container {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  button {
    flex-shrink: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 2rem;
    width: 2rem;
    border: var(--wa-border-style) var(--wa-border-width-s)
      var(--wa-color-neutral-border-loud);
    border-radius: var(--wa-border-radius-m);
    background: none;
    cursor: pointer;
    font-size: var(--wa-font-size-s);
    font-weight: var(--wa-font-weight-bold);
    color: var(--wa-color-text-quiet);

    &[data-variant] {
      background-color: var(--wa-color-gray-95);

      & wa-spinner {
        --track-color: var(--wa-color-gray-50);
        --indicator-color: var(--wa-color-gray-90);
      }

      border-color: var(--wa-color-gray-50);
      color: var(--wa-color-gray-50);
    }

    &[data-variant="top"] {
      background-color: var(--wa-color-green-95);

      & wa-spinner {
        --track-color: var(--wa-color-green-50);
        --indicator-color: var(--wa-color-green-90);
      }

      border-color: var(--wa-color-green-50);
      color: var(--wa-color-green-50);
    }

    &[data-variant="flash"] {
      background-color: var(--wa-color-yellow-95);

      & wa-spinner {
        --track-color: var(--wa-color-yellow-50);
        --indicator-color: var(--wa-color-yellow-90);
      }

      border-color: var(--wa-color-yellow-50);
      color: var(--wa-color-yellow-50);
    }
  }

  .label {
    display: flex;
    align-items: center;
    gap: var(--wa-space-s);
  }

  button:disabled {
    cursor: not-allowed;
    border: 0;
  }

  wa-dialog {
    &::part(body) {
      display: flex;
      flex-direction: column;
      gap: var(--wa-space-m);
    }

    & .horizontal {
      display: flex;
      gap: var(--wa-space-s);
    }
  }

  .horizontal {
    width: 100%;
    display: flex;
    gap: var(--wa-space-2xs);

    & wa-button {
      flex: 1;
    }
  }

  .small {
    font-size: var(--wa-font-size-xs);
  }
</style>
