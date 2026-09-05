<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import { HoldColorIndicator } from "@climblive/lib/components";
  import type { PointValue, Problem, Tick } from "@climblive/lib/models";
  import { deleteTickMutation, putTickMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { AxiosError } from "axios";
  import { getContext } from "svelte";
  import type { Readable } from "svelte/store";
  import SimpleTickEditor from "./tick/SimpleTickEditor.svelte";
  import TickEditor from "./tick/TickEditor.svelte";

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
  const putTick = $derived(putTickMutation($session.contenderId, problem.id));
  const deleteTick = $derived(deleteTickMutation());

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
    if (!tick?.id) {
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
      <TickEditor {problem} {tick} {pointValue} {showPoints} {putTick} />
    {:else}
      <SimpleTickEditor {problem} {tick} {pointValue} {putTick} bind:open />
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
        disabled={open && tick?.id === undefined}
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
</style>
