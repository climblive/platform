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

  let loading = $derived($createTick.isPending || $deleteTick.isPending);
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
      $deleteTick.mutate(tick.id, {
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

  const handleTick = (event: MouseEvent, flash: boolean) => {
    event.stopPropagation();

    navigator.vibrate?.(50);
    open = false;

    $createTick.mutate(
      {
        problemId: problem.id,
        top: true,
        attemptsTop: flash ? 1 : 999,
        zone: true,
        attemptsZone: flash ? 1 : 999,
      },
      {
        onError: (error) => {
          if (error instanceof AxiosError && error.status === 409) {
            toastError("Ascent is already registered.");
          } else {
            toastError("Failed to register ascent.");
          }
        },
      },
    );
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
      <wa-icon name="lightning-charge"></wa-icon>
    {:else if variant === "ticked"}
      <wa-icon name="check2-all"></wa-icon>
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
    <wa-button size="small" onclick={(e: MouseEvent) => handleTick(e, false)}>
      <wa-icon slot="prefix" name="check2-all"></wa-icon>
      Top
    </wa-button>
    <wa-button size="small" onclick={(e: MouseEvent) => handleTick(e, true)}>
      <wa-icon slot="prefix" name="lightning-charge"></wa-icon>
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
    height: calc(100% - 2 * var(--wa-spacing-2x-small));
    aspect-ratio: 1 / 1;
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
    border-color: var(--wa-color-gray-600);
    border-width: 2px;
    border-style: dotted;
    border-radius: var(--wa-border-radius-small);
  }

  button:disabled {
    cursor: not-allowed;
    border: 0;
  }

  div[data-variant] {
    & > button {
      border-width: calc(2 * var(--wa-input-border-width));
    }
  }

  div[data-variant="ticked"] {
    background-color: var(--wa-color-green-100);

    & wa-spinner {
      --track-color: var(--wa-color-green-600);
      --indicator-color: var(--wa-color-green-100);
    }

    & > button {
      border-color: var(--wa-color-green-600);
      color: var(--wa-color-green-600);
    }
  }

  div[data-variant="flashed"] {
    background-color: var(--wa-color-yellow-100);

    & wa-spinner {
      --track-color: var(--wa-color-yellow-500);
      --indicator-color: var(--wa-color-yellow-100);
    }

    & > button {
      border-color: var(--wa-color-yellow-500);
      color: var(--wa-color-yellow-500);
    }
  }

  wa-popup {
    --arrow-color: var(--wa-color-primary-600);

    & wa-button::part(base) {
      width: 2.5rem;
      height: 2.5rem;
      flex-direction: column;
      align-items: center;
      padding: 0;
    }

    & wa-button::part(prefix) {
      font-size: var(--wa-font-size-small);
    }

    & wa-button::part(label) {
      font-size: var(--wa-font-size-2x-small);
      line-height: var(--wa-line-height-dense);
      padding: 0;
    }
  }

  wa-popup::part(popup) {
    background-color: var(--wa-color-primary-600);
    box-shadow:
      rgba(50, 50, 93, 0.25) 0px 13px 27px -5px,
      rgba(0, 0, 0, 0.3) 0px 8px 16px -8px;
    border-radius: var(--wa-border-radius-medium);
    padding: var(--wa-spacing-small);
    cursor: default;
  }

  wa-popup[active]::part(popup) {
    display: flex;
    gap: var(--wa-spacing-2x-small);
  }
</style>
