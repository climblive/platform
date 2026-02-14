<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
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
    required = false,
    allowClear = false,
    name,
    ...rest
  }: Props = $props();

  const id = $props.id();
  let value = $state(rest.value);

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
    if (!hiddenInput) {
      return;
    }

    value = color;
    hiddenInput.value = color;
    popover?.hide();

    if (hiddenInput) {
      hiddenInput.dispatchEvent(new Event("input", { bubbles: true }));
    }
  };

  const handleClear = () => {
    if (!hiddenInput) {
      return;
    }

    value = undefined;
    hiddenInput.value = "";
    popover?.hide();

    if (hiddenInput) {
      hiddenInput.dispatchEvent(new Event("input", { bubbles: true }));
    }
  };
</script>

<div class="hold-color-picker">
  <label for={id}>{label}</label>
  <input bind:this={hiddenInput} type="hidden" {name} {required} {value} />

  <button {id} class="trigger-button" type="button">
    <HoldColorIndicator --height="2rem" --width="2rem" primary={value} />
  </button>

  <wa-popover bind:this={popover} for={id} placement="right" distance={10}>
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
            <HoldColorIndicator
              --height="1.25rem"
              --width="1.25rem"
              primary={color}
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
    gap: var(--wa-space-2xs);
  }

  label {
    font-size: var(--wa-font-size-s);
    font-weight: var(--wa-form-control-label-font-weight);
    color: var(--wa-form-control-label-color);
    line-height: var(--wa-form-control-label-line-height);
  }

  .trigger-button {
    border: none;
    background: transparent;
    width: max-content;
    cursor: pointer;
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
  }

  .color-button {
    background-color: transparent;
    cursor: pointer;
    border: none;
    padding: var(--wa-space-2xs);
  }

  .clear-button {
    margin-block-start: var(--wa-space-s);
    width: 100%;
  }
</style>
