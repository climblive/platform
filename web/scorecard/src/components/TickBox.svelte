<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import type WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
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

  let dialog: WaDialog | undefined = $state();

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
      case tick?.zone2:
        return "zone2";
      case tick?.zone1:
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
    event: InputEvent,
    feature: "zone1" | "zone2" | "top",
    flash: boolean,
  ) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    const uncheck = !(event.target as WaCheckbox).checked;

    if (uncheck) {
      tickBuilder.unreachFeature(feature);
    } else {
      tickBuilder.reachFeature(feature);
    }

    const nextTick = tickBuilder.tick;

    if (!enableAttempts) {
      let attempts = 999;
      if (flash && !uncheck) {
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
    {:else if variant === "flash"}
      <pre>F</pre>
    {:else if variant === "top"}
      <pre>T</pre>
    {:else if variant === "zone2"}
      <pre>Z2</pre>
    {:else if variant === "zone1"}
      <pre>Z1</pre>
    {/if}
  </button>

  <wa-dialog
    label="Problem number {problem.number}"
    bind:this={dialog}
    {open}
    light-dismiss
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
        onChange={(e: InputEvent) => handleTick(e, "top", false)}
        points={showPoints ? pointValue?.top : undefined}
        checked={tick?.top}
        indeterminate={tick?.top && tick?.attemptsTop === 1}
        attempts={enableAttempts ? (tick?.attemptsTop ?? 0) : undefined}
      />

      {#if !enableAttempts}
        <TickButton
          label="Flash"
          onChange={(e: InputEvent) => handleTick(e, "top", true)}
          points={showPoints ? pointValue?.top : undefined}
          bonusPoints={pointValue?.flashBonus}
          checked={tick?.top && tick?.attemptsTop === 1}
          attempts={enableAttempts ? (tick?.attemptsTop ?? 0) : undefined}
        />
      {/if}
    </div>

    {#if problem.zone2Enabled}
      <TickButton
        label="Zone 2"
        onChange={(e: InputEvent) => handleTick(e, "zone2", false)}
        points={showPoints ? pointValue?.zone2 : undefined}
        checked={tick?.zone2}
        indeterminate={tick?.top}
        attempts={enableAttempts ? (tick?.attemptsZone2 ?? 0) : undefined}
      />
    {/if}

    {#if problem.zone1Enabled}
      <TickButton
        label="Zone 1"
        onChange={(e: InputEvent) => handleTick(e, "zone1", false)}
        points={showPoints ? pointValue?.zone1 : undefined}
        checked={tick?.zone1}
        indeterminate={tick?.zone2}
        attempts={enableAttempts ? (tick?.attemptsZone1 ?? 0) : undefined}
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
    display: flex;
    justify-content: center;
    align-items: center;
    height: calc(100% - 2 * var(--wa-space-xs));
    aspect-ratio: 1 / 1;
    border: var(--wa-border-style) var(--wa-border-width-s)
      var(--wa-color-neutral-border-loud);
    border-radius: var(--wa-border-radius-l);
    background: none;
    cursor: pointer;
    width: max-content;
    font-size: var(--wa-font-size-s);
    font-weight: var(--wa-font-weight-bold);

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

    & pre {
      margin: 0;
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
</style>
