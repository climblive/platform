<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import type WaPopup from "@awesome.me/webawesome/dist/components/popup/popup.js";
  import type { Problem, Tick } from "@climblive/lib/models";
  import {
    createTickMutation,
    deleteTickMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { AxiosError } from "axios";
  import { getContext } from "svelte";
  import type { Readable } from "svelte/store";

  interface Props {
    problem: Problem;
    tick: Tick | undefined;
    disabled: boolean | undefined;
  }

  let { problem, tick, disabled = false }: Props = $props();

  let container: HTMLDivElement | undefined = $state();
  let popup: WaPopup | undefined = $state();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");
  const createTick = $derived(createTickMutation($session.contenderId));
  const deleteTick = $derived(deleteTickMutation());

  let open = $state(false);

  let loading = $derived(createTick.isPending || deleteTick.isPending);
  let variant = $derived(
    tick ? (tick.attemptsTop === 1 ? "flashed" : "ticked") : undefined,
  );

  const handleClickOutside = (event: MouseEvent) => {
    if (
      popup &&
      event.target instanceof Node &&
      !popup.contains(event.target)
    ) {
      open = false;
    }
  };

  const handleCheck = () => {
    if (tick?.id) {
      deleteTick.mutate(tick.id, {
        onError: (error) => {
          if (error instanceof AxiosError && error.status === 404) {
            toastError("Ascent is already removed.");
          } else {
            toastError("Failed to remove ascent.");
          }
        },
      });
    } else {
      open = true;
    }
  };

  const handleTick = (
    event: MouseEvent,
    feature: "zone1" | "zone2" | "top",
    flash: boolean,
  ) => {
    event.stopPropagation();

    navigator.vibrate?.(50);
    open = false;

    const tick: Omit<Tick, "id" | "timestamp"> = {
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
        tick.top = true;
        tick.zone2 = true;
        tick.zone1 = true;
        break;
      case "zone2":
        tick.zone2 = true;
        tick.zone1 = true;
      case "zone1":
        tick.zone1 = true;
    }

    createTick.mutate(tick, {
      onError: (error) => {
        if (error instanceof AxiosError && error.status === 409) {
          toastError("Ascent is already registered.");
        } else {
          toastError("Failed to register ascent.");
        }
      },
    });
  };

  $effect(() => {
    if (popup && container) {
      popup.anchor = container;
    }
  });

  $effect(() => {
    if (tick !== undefined) {
      open = false;
    }
  });
</script>

<svelte:body on:click|capture={handleClickOutside} />

<div data-variant={variant} bind:this={container}>
  <button
    disabled={disabled || loading}
    onclick={handleCheck}
    aria-label={tick?.id ? "Untick" : "Tick"}
  >
    {#if loading}
      <wa-spinner></wa-spinner>
    {:else if variant === "flashed"}
      <wa-icon name="bolt"></wa-icon>
    {:else if variant === "ticked"}
      <wa-icon name="check"></wa-icon>
    {/if}
  </button>

  <wa-popup
    bind:this={popup}
    placement="left"
    active={open}
    arrow
    strategy="fixed"
    distance="10"
  >
    <wa-button
      size="small"
      onclick={(e: MouseEvent) => handleTick(e, "zone1", false)}
    >
      <wa-icon slot="start" name="check"></wa-icon>
      Zone 1
    </wa-button>
    <wa-button
      size="small"
      onclick={(e: MouseEvent) => handleTick(e, "zone2", false)}
    >
      <wa-icon slot="start" name="check"></wa-icon>
      Zone 2
    </wa-button>
    <wa-button
      size="small"
      onclick={(e: MouseEvent) => handleTick(e, "top", false)}
    >
      <wa-icon slot="start" name="check"></wa-icon>
      Top
    </wa-button>
    <wa-button
      size="small"
      onclick={(e: MouseEvent) => handleTick(e, "top", true)}
    >
      <wa-icon slot="start" name="bolt"></wa-icon>
      Flash
    </wa-button>
  </wa-popup>
</div>

<style>
  div {
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
    height: calc(100% - 2 * var(--wa-space-2xs));
    aspect-ratio: 1 / 1;
    border-radius: var(--wa-border-radius-s);
  }

  button {
    box-sizing: content-box;
    border: none;
    background: none;
    padding: 0;
    width: 1.5rem;
    height: 1.5rem;
    cursor: pointer;
    display: flex;
    justify-content: center;
    align-items: center;
    border-color: var(--wa-color-neutral-border-loud);
    border-width: var(--wa-border-width-s);
    border-style: var(--wa-border-style);
    border-radius: var(--wa-border-radius-s);
  }

  button:disabled {
    cursor: not-allowed;
    border: 0;
  }

  div[data-variant="ticked"] {
    background-color: var(--wa-color-green-95);

    & wa-spinner {
      --track-color: var(--wa-color-green-50);
      --indicator-color: var(--wa-color-green-90);
    }

    & > button {
      border-color: var(--wa-color-green-50);
      color: var(--wa-color-green-50);
    }
  }

  div[data-variant="flashed"] {
    background-color: var(--wa-color-yellow-95);

    & wa-spinner {
      --track-color: var(--wa-color-yellow-50);
      --indicator-color: var(--wa-color-yellow-90);
    }

    & > button {
      border-color: var(--wa-color-yellow-50);
      color: var(--wa-color-yellow-50);
    }
  }

  wa-popup {
    --arrow-color: var(--wa-color-brand-fill-loud);

    & wa-button::part(base) {
      width: 2.5rem;
      height: 2.5rem;
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 0;
    }

    & wa-button > wa-icon {
      margin: 0;
    }

    & wa-button::part(label) {
      font-size: var(--wa-font-size-2xs);
      line-height: var(--wa-line-height-condensed);
    }
  }

  wa-popup::part(popup) {
    background-color: var(--wa-color-brand-fill-loud);
    box-shadow: var(--wa-shadow-l);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-s);
    cursor: default;
  }

  wa-popup[active]::part(popup) {
    display: flex;
    gap: var(--wa-space-2xs);
  }
</style>
