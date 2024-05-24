<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import type { Problem, Tick } from "@climblive/lib/models";
  import {
    createTickMutation,
    deleteTickMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { SlPopup } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/popup/popup.js";
  import "@shoelace-style/shoelace/dist/components/spinner/spinner.js";
  import { afterUpdate, getContext } from "svelte";
  import type { Readable } from "svelte/store";

  export let problem: Problem;
  export let tick: Tick | undefined;

  let container: HTMLDivElement;
  let popup: SlPopup;

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");
  const createTick = createTickMutation($session.contenderId);
  const deleteTick = deleteTickMutation();

  let open = false;

  $: loading = $createTick.isPending || $deleteTick.isPending;
  $: variant = !!tick ? (tick.flash ? "flashed" : "ticked") : undefined;

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
        onError: () => toastError("Failed to remove ascent."),
      });
    } else {
      open = true;
    }
  };

  const handleTick = (flash: boolean) => {
    navigator.vibrate(50);
    open = false;

    $createTick.mutate(
      {
        problemId: problem.id,
        contenderId: $session.contenderId,
        flash,
      },
      {
        onError: () => toastError("Failed to register ascent."),
      }
    );
  };

  afterUpdate(() => {
    popup.anchor = container;
  });
</script>

<svelte:body on:click|capture={handleClickOutside} />

<div data-variant={variant} bind:this={container}>
  <button disabled={loading} on:click={handleCheck}>
    {#if loading}
      <sl-spinner></sl-spinner>
    {:else if variant === "flashed"}
      <sl-icon name="lightning-charge"></sl-icon>
    {:else if variant === "ticked"}
      <sl-icon name="check2-all"></sl-icon>
    {/if}
  </button>

  <sl-popup
    bind:this={popup}
    placement="left"
    active={open}
    arrow
    strategy="fixed"
    distance="4"
  >
    <sl-button size="small" on:click|stopPropagation={() => handleTick(true)}>
      <sl-icon slot="prefix" name="lightning-charge"></sl-icon>
      Flash
    </sl-button>
    <sl-button size="small" on:click|stopPropagation={() => handleTick(false)}>
      <sl-icon slot="prefix" name="check2-all"></sl-icon>
      Top
    </sl-button>
  </sl-popup>
</div>

<style>
  div {
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
    height: calc(100% - 2 * var(--sl-spacing-2x-small));
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
    border-color: var(--sl-color-primary-900);
    border-width: var(--sl-input-border-width);
    border-style: solid;
    border-radius: var(--sl-border-radius-small);
  }

  button:disabled {
    cursor: not-allowed;
  }

  div[data-variant] {
    & > button {
      border-width: calc(2 * var(--sl-input-border-width));
    }
  }

  div[data-variant="ticked"] {
    background-color: var(--sl-color-green-100);

    & sl-spinner {
      --track-color: var(--sl-color-green-600);
      --indicator-color: var(--sl-color-green-100);
    }

    & > button {
      border-color: var(--sl-color-green-600);
      color: var(--sl-color-green-600);
    }
  }

  div[data-variant="flashed"] {
    background-color: var(--sl-color-yellow-100);

    & sl-spinner {
      --track-color: var(--sl-color-yellow-500);
      --indicator-color: var(--sl-color-yellow-100);
    }

    & > button {
      border-color: var(--sl-color-yellow-500);
      color: var(--sl-color-yellow-500);
    }
  }

  sl-popup {
    --arrow-color: var(--sl-color-primary-800);

    & sl-button::part(base) {
      width: 2.5rem;
      height: 2.5rem;
      flex-direction: column;
      align-items: center;
      padding: 0;
    }

    & sl-button::part(prefix) {
      font-size: var(--sl-font-size-small);
    }

    & sl-button::part(label) {
      font-size: var(--sl-font-size-2x-small);
      line-height: var(--sl-line-height-dense);
      padding: 0;
    }
  }

  sl-popup::part(popup) {
    background-color: var(--sl-color-primary-800);
    box-shadow:
      rgba(50, 50, 93, 0.25) 0px 13px 27px -5px,
      rgba(0, 0, 0, 0.3) 0px 8px 16px -8px;
    border-radius: var(--sl-border-radius-medium);
    padding: var(--sl-spacing-small);
    cursor: default;
  }

  sl-popup[active]::part(popup) {
    display: flex;
    gap: var(--sl-spacing-2x-small);
  }
</style>
