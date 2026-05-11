<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import type WaPopover from "@awesome.me/webawesome/dist/components/popover/popover.js";
  import HoldColorIndicator from "./HoldColorIndicator.svelte";

  interface Props {
    label: string;
    value: string | undefined;
    required?: boolean;
    allowClear?: boolean;
    name: string;
    placement: "top" | "bottom" | "left" | "right";
  }

  const {
    label,
    value: initialValue,
    required = false,
    allowClear = false,
    name,
    placement,
  }: Props = $props();

  const id = $props.id();

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

  let selectedColor = $derived(initialValue);
  let popover = $state<WaPopover>();
  let hiddenInput = $state<HTMLInputElement>();

  const handleColorSelect = (color: string) => {
    selectedColor = color;

    if (hiddenInput) {
      hiddenInput.value = color;
      hiddenInput.dispatchEvent(new Event("input", { bubbles: true }));
    }

    popover?.hide();
  };

  const handleClear = () => {
    selectedColor = undefined;

    if (hiddenInput) {
      hiddenInput.value = "";
      hiddenInput.dispatchEvent(new Event("input", { bubbles: true }));
    }

    popover?.hide();
  };
</script>

<div class="hold-color-picker">
  <label for={id} class:required>{label}</label>
  <input
    bind:this={hiddenInput}
    type="hidden"
    {name}
    {required}
    value={initialValue}
  />

  <button {id} class="trigger" type="button">
    <HoldColorIndicator
      --height="2rem"
      --width="2rem"
      primary={selectedColor}
      outlined
    />
  </button>

  <wa-popover bind:this={popover} for={id} {placement} distance={10}>
    <div class="popup-content" role="menu" aria-label="Color selection">
      <div class="color-grid">
        {#each colors as color (color)}
          <button
            type="button"
            class:selected={selectedColor === color}
            onclick={() => handleColorSelect(color)}
            aria-label="Select color {color}"
            role="menuitem"
          >
            <HoldColorIndicator
              --height="2rem"
              --width="2rem"
              primary={color}
              outlined
            />
          </button>
        {/each}
      </div>
      {#if allowClear && !required}
        <wa-button
          class="clear-button"
          size="small"
          appearance="outlined"
          onclick={handleClear}
          aria-label="Clear color selection"
        >
          Clear
          <wa-icon slot="start" name="xmark"></wa-icon>
        </wa-button>
      {/if}
    </div>
  </wa-popover>
</div>

<style>
  .hold-color-picker {
    display: flex;
    flex-direction: column;
    justify-content: start;
    gap: var(--wa-space-xs);
  }

  label {
    font-size: var(--wa-font-size-s);
    font-weight: var(--wa-form-control-label-font-weight);
    color: var(--wa-form-control-label-color);
    line-height: var(--wa-form-control-label-line-height);
  }

  label.required::after {
    content: var(--wa-form-control-required-content);
    color: var(--wa-form-control-required-content-color);
    margin-inline-start: var(--wa-form-control-required-content-offset);
  }

  .trigger {
    border: none;
    background: transparent;
    width: max-content;
    cursor: pointer;
    padding: 0;
    transition: opacity var(--wa-transition-fast);
  }

  wa-popover::part(body) {
    background-color: var(--wa-color-surface-raised);
    border: var(--wa-border-style) var(--wa-border-width-s)
      var(--wa-color-surface-border);
    padding: var(--wa-space-s);
  }

  wa-popover::part(popup__arrow) {
    background-color: var(--wa-color-surface-raised);
    border-color: var(--wa-color-surface-border);
  }

  .color-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: var(--wa-space-xs);

    button {
      background-color: transparent;
      cursor: pointer;
      border: none;
      padding: 0;
    }
  }

  .clear-button {
    margin-block-start: var(--wa-space-s);
    width: 100%;
  }
</style>
