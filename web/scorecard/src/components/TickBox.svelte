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
  import { TickBuilder } from "./tick.svelte";

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
    disabled = false,
    pointValue,
    showPoints,
    enableAttempts,
    ...rest
  }: Props = $props();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");
  const putTick = $derived(putTickMutation($session.contenderId));
  const deleteTick = $derived(deleteTickMutation());

  const tickId = $derived(rest.tick?.id);
  const tickBuilder = $derived(new TickBuilder(problem, rest.tick));
  const tick = $derived(tickBuilder.tick);

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
    if (!tickId) {
      return;
    }

    event.stopPropagation();

    open = false;

    deleteTick.mutate(tickId, {
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

  const renderAttempts = (attempts: number, showFlash: boolean) => {
    switch (true) {
      case attempts === 1 && showFlash:
        return "Flash";
      case attempts === 1:
        return "1 attempt";
      default:
        return `${attempts} attempts`;
    }
  };
</script>

<div class="container">
  <button
    data-variant={variant}
    disabled={disabled || loading}
    onclick={() => (open = true)}
    aria-label={tickId ? "Edit" : "Tick"}
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
    {:else if tick !== undefined && enableAttempts && tick.attemptsTop}
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
      /> Problem #{problem.number}
    </div>
    {#if enableAttempts}
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
          failed {tick.attemptsTop === 1 ? "attempt" : "attempts"}
        </div>

        <wa-button
          size="m"
          pill
          appearance="outlined"
          onclick={(event: MouseEvent) => handleAddAttempt(event)}
          disabled={!tickBuilder.canAddAttempt()}
        >
          <wa-icon name="plus"></wa-icon>
        </wa-button>
      </div>
    {/if}

    <TickButton
      label={`Top (${renderAttempts((tick?.attemptsTop ?? 0) + 1, true)})`}
      subLabel="T"
      onChange={(checked, flash) => handleTick(checked, "top", flash)}
      points={showPoints ? pointValue?.top : undefined}
      bonusPoints={pointValue?.flashBonus}
      checked={tick?.top}
      attempts={tick?.attemptsTop ?? 0}
      flashToggle={!enableAttempts}
      showAttempts={enableAttempts}
      {showPoints}
    />

    {#if problem.zone2Enabled}
      <TickButton
        label={`Zone 2 (${renderAttempts((tick?.attemptsZone2 ?? 0) + 1, false)})`}
        subLabel="Z₂"
        onChange={(checked) => handleTick(checked, "zone2", false)}
        points={showPoints ? pointValue?.zone2 : undefined}
        checked={tick?.zone2}
        attempts={tick?.attemptsZone2 ?? 0}
        showAttempts={enableAttempts}
        {showPoints}
        disabled={tick?.top}
      />
    {/if}

    {#if problem.zone1Enabled}
      <TickButton
        label={`Zone 1 (${renderAttempts((tick?.attemptsZone1 ?? 0) + 1, false)})`}
        subLabel="Z₁"
        onChange={(checked) => handleTick(checked, "zone1", false)}
        points={showPoints ? pointValue?.zone1 : undefined}
        checked={tick?.zone1}
        attempts={tick?.attemptsZone1 ?? 0}
        showAttempts={enableAttempts}
        {showPoints}
        disabled={tick?.zone2}
      />
    {/if}

    <div class="footer">
      <div class="status">
        {#if putTick.isPending}
          <wa-spinner></wa-spinner>
        {:else if putTick.isSuccess}
          <div class="success">
            <wa-icon name="check"></wa-icon> Saved
          </div>
        {/if}
      </div>

      <wa-button
        size="s"
        appearance="plain"
        onclick={(e: MouseEvent) => handleDelete(e)}
        variant="danger"
        disabled={open && tickId === undefined}
      >
        <wa-icon name="trash" label="Remove"></wa-icon>
      </wa-button>
    </div>
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
  }

  .horizontal {
    margin-inline: auto;
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
  }

  .small {
    font-size: var(--wa-font-size-xs);
  }

  .footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;

    .success {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: var(--wa-space-2xs);
      color: var(--wa-color-success-fill-loud);
      animation: hide 0s ease 2s;
      animation-fill-mode: forwards;
    }
  }

  @keyframes hide {
    to {
      visibility: hidden;
    }
  }

  .attempts {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: var(--wa-font-size-s);
    width: 4rem;
    color: var(--wa-color-text-quiet);

    & .number {
      margin: 0;
      font-size: var(--wa-font-size-l);
      font-weight: var(--wa-font-weight-bold);
      color: var(--wa-color-text-normal);
    }
  }
</style>
