<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { Problem, Tick } from "@climblive/lib/models";
  import { putTickMutation } from "@climblive/lib/queries";
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
  const putTick = $derived(putTickMutation($session.contenderId));
  let open = $state(false);

  const loading = $derived(putTick.isPending);
  const variant = $derived.by((): "top" | "zone2" | "zone1" | undefined => {
    switch (true) {
      case tick?.top:
        return "top";
      case tick?.zone2:
        return "zone2";
      case tick?.zone1:
        return "zone1";
    }
  });
  const scorecardVariant = $derived.by(
    (): "flash" | "top" | "zone2" | "zone1" | undefined => {
      if (tick?.top && tick.attemptsTop === 1) {
        return "flash";
      }

      return variant;
    },
  );
  const attempts = $derived(tick?.attemptsTop ?? 0);
  const hasAttempts = $derived(attempts > 0);
  const displayVariant = $derived.by(
    (): "attempts" | "flash" | "top" | "zone2" | "zone1" | undefined => {
      if (scorecardVariant !== undefined) {
        return scorecardVariant;
      }

      if (hasAttempts) {
        return "attempts";
      }
    },
  );
  const getAttempts = (buttonVariant: "top" | "zone2" | "zone1") => {
    switch (buttonVariant) {
      case "top":
        return tick?.attemptsTop ?? 0;
      case "zone2":
        return tick?.attemptsZone2 ?? 0;
      case "zone1":
        return tick?.attemptsZone1 ?? 0;
    }
  };
  const attemptFeatures = ["top", "zone2", "zone1"] as const;
  const getAttemptField = (feature: (typeof attemptFeatures)[number]) => {
    switch (feature) {
      case "top":
        return "attemptsTop";
      case "zone2":
        return "attemptsZone2";
      case "zone1":
        return "attemptsZone1";
    }
  };
  const getMinimumAttempts = (
    nextTick: Omit<Tick, "id" | "timestamp">,
    feature: (typeof attemptFeatures)[number],
  ) => {
    switch (feature) {
      case "top":
        return nextTick.top ? 1 : 0;
      case "zone2":
        return nextTick.zone2 ? 1 : 0;
      case "zone1":
        return nextTick.zone1 ? 1 : 0;
    }
  };
  const isChecked = (buttonVariant: "top" | "zone2" | "zone1") =>
    variant === buttonVariant;
  const isIndeterminate = (buttonVariant: "zone2" | "zone1") => {
    switch (buttonVariant) {
      case "zone2":
        return variant === "top";
      case "zone1":
        return variant === "top" || variant === "zone2";
    }
  };
  const canAddAttempt = $derived(!tick?.top);

  const getNextTick = (): Omit<Tick, "id" | "timestamp"> => ({
    problemId: problem.id,
    top: tick?.top ?? false,
    zone2: tick?.zone2 ?? false,
    zone1: tick?.zone1 ?? false,
    attemptsTop: tick?.attemptsTop ?? 0,
    attemptsZone2: tick?.attemptsZone2 ?? 0,
    attemptsZone1: tick?.attemptsZone1 ?? 0,
  });

  const incrementAttempts = (nextTick: Omit<Tick, "id" | "timestamp">) => {
    if (!nextTick.top) {
      nextTick.attemptsTop += 1;
    }

    if (!nextTick.zone2) {
      nextTick.attemptsZone2 += 1;
    }

    if (!nextTick.zone1) {
      nextTick.attemptsZone1 += 1;
    }
  };

  const syncAttemptCounts = (
    nextTick: Omit<Tick, "id" | "timestamp">,
    feature: (typeof attemptFeatures)[number],
  ) => {
    const featureIndex = attemptFeatures.indexOf(feature);
    const targetField = getAttemptField(feature);

    nextTick[targetField] = Math.max(
      getMinimumAttempts(nextTick, feature),
      nextTick[targetField],
    );

    for (let index = featureIndex - 1; index >= 0; index -= 1) {
      const currentField = getAttemptField(attemptFeatures[index]);
      const belowField = getAttemptField(attemptFeatures[index + 1]);

      nextTick[currentField] = Math.max(
        nextTick[currentField],
        nextTick[belowField],
      );
    }

    for (
      let index = featureIndex + 1;
      index < attemptFeatures.length;
      index += 1
    ) {
      const currentField = getAttemptField(attemptFeatures[index]);
      const aboveField = getAttemptField(attemptFeatures[index - 1]);

      nextTick[currentField] = Math.min(
        nextTick[currentField],
        nextTick[aboveField],
      );
    }

    for (const currentFeature of attemptFeatures) {
      const currentField = getAttemptField(currentFeature);

      nextTick[currentField] = Math.max(
        getMinimumAttempts(nextTick, currentFeature),
        nextTick[currentField],
      );
    }
  };

  const putNextTick = (nextTick: Omit<Tick, "id" | "timestamp">) => {
    putTick.mutate(nextTick, {
      onError: (error) => {
        if (error instanceof AxiosError && error.status === 409) {
          toastError("Ascent is already registered.");
        } else {
          toastError("Failed to save ascent.");
        }
      },
    });
  };

  const handleAttempt = (event: MouseEvent) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    const nextTick = getNextTick();

    incrementAttempts(nextTick);
    putNextTick(nextTick);
  };

  const handleSetAttempts = (
    event: MouseEvent,
    feature: (typeof attemptFeatures)[number],
    attempts: number,
  ) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    const nextTick = getNextTick();
    const attemptField = getAttemptField(feature);
    const nextAttempts = Number.isFinite(attempts)
      ? Math.trunc(attempts)
      : nextTick[attemptField];

    nextTick[attemptField] = Math.max(
      getMinimumAttempts(nextTick, feature),
      nextAttempts,
    );

    syncAttemptCounts(nextTick, feature);
    putNextTick(nextTick);
  };

  const handleTick = (
    event: MouseEvent,
    feature: "zone1" | "zone2" | "top",
  ) => {
    event.stopPropagation();

    navigator.vibrate?.(50);

    const nextTick = getNextTick();

    if (isChecked(feature)) {
      switch (feature) {
        case "top":
          nextTick.top = false;
          break;
        case "zone2":
          nextTick.zone2 = false;
          break;
        case "zone1":
          nextTick.zone1 = false;
      }

      putNextTick(nextTick);
      return;
    }

    incrementAttempts(nextTick);

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

    putNextTick(nextTick);
  };
</script>

<div class="container">
  <button
    data-variant={displayVariant}
    disabled={disabled || loading}
    onclick={() => (open = true)}
    aria-label={tick?.id ? "Edit" : "Tick"}
  >
    {#if loading}
      <wa-spinner></wa-spinner>
    {:else if scorecardVariant === "flash"}
      <pre>F</pre>
    {:else if scorecardVariant === "top"}
      <pre>T</pre>
    {:else if scorecardVariant === "zone2"}
      <pre>Z2</pre>
    {:else if scorecardVariant === "zone1"}
      <pre>Z1</pre>
    {:else if hasAttempts}
      <pre>{attempts}</pre>
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

    <TickButton
      label="Top"
      onClick={(e: MouseEvent) => handleTick(e, "top")}
      onSaveAttempts={(e: MouseEvent, attempts: number) =>
        handleSetAttempts(e, "top", attempts)}
      points={problem.pointsTop}
      attempts={getAttempts("top")}
      minimumAttempts={getMinimumAttempts(getNextTick(), "top")}
      checked={isChecked("top")}
    />

    {#if problem.zone2Enabled}
      <TickButton
        label="Zone 2"
        onClick={(e: MouseEvent) => handleTick(e, "zone2")}
        onSaveAttempts={(e: MouseEvent, attempts: number) =>
          handleSetAttempts(e, "zone2", attempts)}
        points={problem.pointsZone2}
        attempts={getAttempts("zone2")}
        minimumAttempts={getMinimumAttempts(getNextTick(), "zone2")}
        checked={isChecked("zone2")}
        indeterminate={isIndeterminate("zone2")}
        disabled={isIndeterminate("zone2")}
      />
    {/if}

    {#if problem.zone1Enabled}
      <TickButton
        label="Zone 1"
        onClick={(e: MouseEvent) => handleTick(e, "zone1")}
        onSaveAttempts={(e: MouseEvent, attempts: number) =>
          handleSetAttempts(e, "zone1", attempts)}
        points={problem.pointsZone1}
        attempts={getAttempts("zone1")}
        minimumAttempts={getMinimumAttempts(getNextTick(), "zone1")}
        checked={isChecked("zone1")}
        indeterminate={isIndeterminate("zone1")}
        disabled={isIndeterminate("zone1")}
      />
    {/if}

    {#if open}
      <wa-button
        size="s"
        appearance="outlined"
        onclick={(e: MouseEvent) => handleAttempt(e)}
        disabled={!canAddAttempt}
      >
        Log failed attempt
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
      gap: var(--wa-space-s);
    }
  }
</style>
