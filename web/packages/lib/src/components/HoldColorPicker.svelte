<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import type WaPopover from "@awesome.me/webawesome/dist/components/popover/popover.js";
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

  let popover: WaPopover | undefined = $state();
  let hiddenInput: HTMLInputElement | undefined = $state();

  const handleColorSelect = (color: string) => {
    value = color;
    popover?.hide();
    if (hiddenInput) {
      hiddenInput.dispatchEvent(new Event("input", { bubbles: true }));
    }
  };

  const handleClear = () => {
    value = undefined;
    popover?.hide();
    if (hiddenInput) {
      hiddenInput.dispatchEvent(new Event("input", { bubbles: true }));
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
    id="{name}-trigger"
    size="small"
    appearance="plain"
    aria-labelledby="{name}-label"
  >
    <div class="trigger-content">
      {#if value}
        <HoldColorIndicator primary={value} />
      {:else}
        <div class="placeholder">
          <span>Select color</span>
        </div>
      {/if}
    </div>
  </wa-button>

  <wa-popover
    bind:this={popover}
    for="{name}-trigger"
    placement="top-start"
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
            appearance="outlined"
            onclick={handleClear}
            style="width: 100%;"
            aria-label="Clear color selection"
          >
            Clear
          </wa-button>
        </div>
      {/if}
    </div>
  </wa-popover>
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
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
  }

  .placeholder::before {
    content: "";
    width: 24px;
    height: 24px;
    background-image: linear-gradient(45deg, #ccc 25%, transparent 25%),
      linear-gradient(-45deg, #ccc 25%, transparent 25%),
      linear-gradient(45deg, transparent 75%, #ccc 75%),
      linear-gradient(-45deg, transparent 75%, #ccc 75%);
    background-size: 8px 8px;
    background-position:
      0 0,
      0 4px,
      4px -4px,
      -4px 0;
    border: 1px solid var(--wa-color-border-normal);
    border-radius: var(--wa-radius-xs);
  }

  .popup-content {
    background: var(--wa-color-surface-border);
  }

  .color-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: var(--wa-space-xs);
    padding: var(--wa-space-s);
  }

  .color-button {
    border: 2px solid transparent;
    border-radius: var(--wa-radius-s);
    background-color: transparent;
    padding: var(--wa-space-3xs);
    cursor: pointer;
    width: 40px;
    height: 40px;
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
    width: 28px;
    height: 28px;
  }

  .clear-section {
    border-top: 1px solid var(--wa-color-border-normal);
    padding: var(--wa-space-s);
    padding-top: var(--wa-space-xs);
    margin-top: var(--wa-space-xs);
  }
</style>
