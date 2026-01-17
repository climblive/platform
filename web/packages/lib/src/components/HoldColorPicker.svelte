<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/popup/popup.js";
  import type WaPopup from "@awesome.me/webawesome/dist/components/popup/popup.js";
  import HoldColorIndicator from "./HoldColorIndicator.svelte";

  interface Props {
    label: string;
    value?: string;
    required?: boolean;
    allowClear?: boolean;
    name?: string;
  }

  let {
    label,
    value = $bindable(),
    required = false,
    allowClear = false,
    name,
  }: Props = $props();

  const colors = [
    "#6f3601",
    "#dc3146",
    "#f46a45",
    "#fac22b",
    "#00ac49",
    "#2fbedc",
    "#0071ec",
    "#9951db",
    "#e66ba3",
    "#9194a2",
    "#000",
    "#fff",
  ];

  let popup: WaPopup | undefined = $state();
  let triggerButton: HTMLElement | undefined = $state();
  let hiddenInput: HTMLInputElement | undefined = $state();

  const handleColorSelect = (color: string) => {
    value = color;
    if (popup) {
      popup.active = false;
    }
  };

  const handleClear = () => {
    value = undefined;
    if (popup) {
      popup.active = false;
    }
  };

  const handleTriggerClick = () => {
    if (popup) {
      popup.active = !popup.active;
    }
  };
</script>

<div class="hold-color-picker">
  <label id="{name}-label">{label}</label>
  <input
    bind:this={hiddenInput}
    type="hidden"
    id={name}
    {name}
    {required}
    value={value ?? ""}
  />
  <wa-button
    size="small"
    appearance="plain"
    onclick={handleTriggerClick}
    bind:this={triggerButton}
    aria-labelledby="{name}-label"
    aria-haspopup="true"
    aria-expanded={popup?.active ?? false}
  >
    <div class="trigger-content">
      {#if value}
        <HoldColorIndicator primary={value} />
      {:else}
        <span class="placeholder">Select color</span>
      {/if}
    </div>
  </wa-button>

  <wa-popup
    bind:this={popup}
    anchor={triggerButton}
    placement="bottom-start"
    distance={4}
  >
    <div class="popup-content" role="listbox" aria-label="Color selection">
      <div class="color-grid">
        {#each colors as color (color)}
          <button
            type="button"
            class="color-button"
            class:selected={value === color}
            onclick={() => handleColorSelect(color)}
            aria-label="Select color {color}"
            role="option"
            aria-selected={value === color}
          >
            <HoldColorIndicator primary={color} />
          </button>
        {/each}
      </div>
      {#if allowClear && !required}
        <div class="clear-section">
          <wa-button
            size="small"
            appearance="plain"
            onclick={handleClear}
            style="width: 100%;"
            aria-label="Clear color selection"
          >
            Clear
          </wa-button>
        </div>
      {/if}
    </div>
  </wa-popup>
</div>

<style>
  .hold-color-picker {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-2xs);
  }

  label {
    font-size: var(--wa-font-size-sm);
    font-weight: var(--wa-font-weight-medium);
    color: var(--wa-color-text-normal);
  }

  .trigger-content {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
    min-width: 120px;
  }

  .trigger-content :global(svg) {
    width: 24px;
    height: 24px;
  }

  .placeholder {
    color: var(--wa-color-text-muted);
  }

  .popup-content {
    background: var(--wa-color-bg-surface);
    border: 1px solid var(--wa-color-border-normal);
    border-radius: var(--wa-radius-m);
    padding: var(--wa-space-s);
    box-shadow: var(--wa-shadow-l);
  }

  .color-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: var(--wa-space-xs);
    margin-bottom: var(--wa-space-xs);
  }

  .color-button {
    background: none;
    border: 2px solid transparent;
    border-radius: var(--wa-radius-s);
    padding: var(--wa-space-2xs);
    cursor: pointer;
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .color-button:hover {
    border-color: var(--wa-color-border-hover);
    background: var(--wa-color-bg-hover);
  }

  .color-button.selected {
    border-color: var(--wa-color-border-active);
    background: var(--wa-color-bg-active);
  }

  .color-button :global(svg) {
    width: 32px;
    height: 32px;
  }

  .clear-section {
    border-top: 1px solid var(--wa-color-border-normal);
    padding-top: var(--wa-space-xs);
  }
</style>
