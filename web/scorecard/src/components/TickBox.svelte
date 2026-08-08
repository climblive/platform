<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { PointValue, Problem, Tick } from "@climblive/lib/models";
  import { deleteTickMutation, putTickMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { AxiosError } from "axios";
  import { getContext } from "svelte";
  import type { Readable } from "svelte/store";
  import TickButton from "./TickButton.svelte";

  interface Props {
    problem: Problem;
    tick: Tick | undefined;
    disabled: boolean | undefined;
    pointValue?: PointValue;
    enableAttempts: boolean;
  }

  let {
    problem,
    tick,
    disabled = false,
    pointValue,
    enableAttempts,
  }: Props = $props();

  let dialog: WaDialog | undefined = $state();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");
  const putTick = $derived(putTickMutation($session.contenderId));
  const deleteTick = $derived(deleteTickMutation());

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
          toastUnexpectedError("Failed to remove ascent.");
        }
      },
    });
  };

  const handleTick = (
    event: MouseEvent,
    feature: "zone1" | "zone2" | "top",
    flash: boolean,
  ) => {
    event.stopPropagation();

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
  };

  const handleLogAttempt = (event: MouseEvent) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    const nextTick: Omit<Tick, "id" | "timestamp"> = {
      problemId: problem.id,
      top: tick?.top ?? false,
      zone2: tick?.zone2 ?? false,
      zone1: tick?.zone1 ?? false,
      attemptsTop: (tick?.attemptsTop ?? 0) + 1,
      attemptsZone2: (tick?.attemptsZone2 ?? 0) + 1,
      attemptsZone1: (tick?.attemptsZone1 ?? 0) + 1,
    };

    putTick.mutate(nextTick, {
      onError: () => {
        toastUnexpectedError("Failed to register ascent.");
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
        onClick={(e: MouseEvent) => handleTick(e, "top", false)}
        points={pointValue?.top}
        active={variant === "top"}
        attempts={enableAttempts ? tick?.attemptsTop : undefined}
      />

      {#if !enableAttempts}
        <TickButton
          label="Flash"
          onClick={(e: MouseEvent) => handleTick(e, "top", true)}
          points={pointValue?.top}
          bonusPoints={pointValue?.flashBonus}
          active={variant === "flash"}
          attempts={enableAttempts ? tick?.attemptsTop : undefined}
        />
      {/if}
    </div>

    {#if problem.zone2Enabled}
      <TickButton
        label="Zone 2"
        onClick={(e: MouseEvent) => handleTick(e, "zone2", false)}
        points={pointValue?.zone2}
        active={variant === "zone2"}
        attempts={enableAttempts ? tick?.attemptsZone2 : undefined}
      />
    {/if}

    {#if problem.zone1Enabled}
      <TickButton
        label="Zone 1"
        onClick={(e: MouseEvent) => handleTick(e, "zone1", false)}
        points={pointValue?.zone1}
        active={variant === "zone1"}
        attempts={enableAttempts ? tick?.attemptsZone1 : undefined}
      />
    {/if}

    {#if enableAttempts}
      <wa-button
        size="s"
        appearance="outlined"
        onclick={(event: MouseEvent) => handleLogAttempt(event)}
      >
        <wa-icon slot="start" name="plus"></wa-icon>
        Attempt
      </wa-button>
    {/if}

    {#if open && variant !== undefined}
      <wa-button
        size="s"
        appearance="plain"
        onclick={(e: MouseEvent) => handleDelete(e)}
        variant="danger"
      >
        <wa-icon slot="start" name="rotate-left"></wa-icon>
        Unsend
      </wa-button>
    {/if}
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
</style>
