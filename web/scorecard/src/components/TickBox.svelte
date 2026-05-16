<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { Problem, Tick, TickPatch } from "@climblive/lib/models";
  import {
    createTickMutation,
    deleteTickMutation,
    updateTickMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { AxiosError } from "axios";
  import { getContext } from "svelte";
  import type { Readable } from "svelte/store";
  import TickButton from "./TickButton.svelte";

  interface Props {
    problem: Problem;
    tick: Tick | undefined;
    disabled: boolean | undefined;
  }

  let { problem, tick, disabled = false }: Props = $props();

  let dialog: WaDialog | undefined = $state();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");
  const createTick = $derived(createTickMutation($session.contenderId));
  const deleteTick = $derived(deleteTickMutation());
  const updateTick = $derived(updateTickMutation($session.contenderId));

  let open = $state(false);

  const loading = $derived(
    createTick.isPending || deleteTick.isPending || updateTick.isPending,
  );
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

  const handleCheck = () => {
    open = true;

    // deleteTick.mutate(tick.id, {
    //   onError: (error) => {
    //     if (error instanceof AxiosError && error.status === 404) {
    //       toastError("Ascent is already removed.");
    //     } else {
    //       toastError("Failed to remove ascent.");
    //     }
    //   },
    // });
  };

  const handleTick = (
    event: MouseEvent,
    feature: "zone1" | "zone2" | "top",
    flash: boolean,
  ) => {
    event.stopPropagation();

    navigator.vibrate?.(50);
    open = false;

    if (tick !== undefined) {
      const patch: TickPatch = {
        top: false,
        zone2: false,
        zone1: false,
        attemptsTop: flash ? 1 : 999,
        attemptsZone2: flash ? 1 : 999,
        attemptsZone1: flash ? 1 : 999,
      };

      switch (feature) {
        case "top":
          patch.top = true;
          patch.zone2 = true;
          patch.zone1 = true;
          break;
        case "zone2":
          patch.zone2 = true;
          patch.zone1 = true;
          break;
        case "zone1":
          patch.zone1 = true;
      }

      updateTick.mutate(
        { tickId: tick.id, patch },
        {
          onError: () => {
            toastError("Failed to register ascent.");
          },
        },
      );
    } else {
      const newTick: Omit<Tick, "id" | "timestamp"> = {
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
          newTick.top = true;
          newTick.zone2 = true;
          newTick.zone1 = true;
          break;
        case "zone2":
          newTick.zone2 = true;
          newTick.zone1 = true;
          break;
        case "zone1":
          newTick.zone1 = true;
      }

      createTick.mutate(newTick, {
        onError: (error) => {
          if (error instanceof AxiosError && error.status === 409) {
            toastError("Ascent is already registered.");
          } else {
            toastError("Failed to register ascent.");
          }
        },
      });
    }
  };

  $effect(() => {
    if (tick !== undefined) {
      open = false;
    }
  });
</script>

<div class="container">
  <button
    data-variant={variant}
    disabled={disabled || loading}
    onclick={handleCheck}
    aria-label={tick?.id ? "Untick" : "Tick"}
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
        iconName="check"
        label="Top"
        onClick={(e: MouseEvent) => handleTick(e, "top", false)}
        points={problem.pointsTop}
        emphasized={variant === "top"}
      />

      <TickButton
        iconName="bolt"
        label="Flash"
        onClick={(e: MouseEvent) => handleTick(e, "top", true)}
        points={problem.pointsTop + (problem.flashBonus ?? 0)}
        emphasized={variant === "flash"}
      />
    </div>

    {#if problem.zone2Enabled}
      <TickButton
        iconName="check"
        label="Zone 2"
        onClick={(e: MouseEvent) => handleTick(e, "zone2", false)}
        points={problem.pointsZone2}
        emphasized={variant === "zone2"}
      />
    {/if}

    {#if problem.zone1Enabled}
      <TickButton
        iconName="check"
        label="Zone 1"
        onClick={(e: MouseEvent) => handleTick(e, "zone1", false)}
        points={problem.pointsZone1}
        emphasized={variant === "zone1"}
      />
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
      gap: var(--wa-space-l);
    }

    & .horizontal {
      display: flex;
      gap: var(--wa-space-s);
    }
  }
</style>
