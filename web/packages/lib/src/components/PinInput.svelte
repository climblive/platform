<script lang="ts">
  export let length: number;
  export let placeholder: string | undefined = undefined;
  export let disabled: boolean = false;
  export let onChange: (value: string) => void;
  export let defaultValue: string | undefined;

  $: inputs = Array.from<HTMLInputElement>({ length });

  const focusInputField = (dir: "next" | "prev", index: number) => {
    if (dir === "next") {
      const nextIndex = index + 1;
      inputs[nextIndex < length ? nextIndex : index].focus();
    }

    if (dir === "prev") {
      const nextIndex = index - 1;
      inputs[nextIndex > -1 ? nextIndex : index].focus();
    }
  };

  const handleInput = (event: InputEvent, index: number) => {
    if (event.isComposing) {
      // Mobile browsers enter composition (IME) when the user starts typing.
      // On Chrome, when focusing the next input field the composition is ended
      // and the input can be continued in the next input field. However, on
      // Firefox for Android, terminating composition also seems to close the
      // keyboard. To force Firefox to end composition before focusing the next
      // input field we blur the current input.
      inputs[index].blur();
    }

    focusInputField("next", index);
  };

  const handleKeyDown = (event: KeyboardEvent, index: number) => {
    if (event.ctrlKey || event.altKey) {
      return;
    }

    if (event.key === "ArrowLeft") {
      event.preventDefault();
      focusInputField("prev", index);
    } else if (event.key === "ArrowRight") {
      event.preventDefault();
      focusInputField("next", index);
    } else if (event.key === "Backspace") {
      event.preventDefault();
      inputs[index].value = "";
      focusInputField("prev", index);
    }
  };

  const handlePaste = (event: ClipboardEvent) => {
    const pasteValue = event.clipboardData?.getData("Text");
    for (const index in inputs) {
      inputs[index].value = pasteValue?.[index] ?? "";
    }
  };

  const handleKeyUp = () => {
    const value = inputs.map((input) => input.value).join("");
    if (value.length === length) {
      onChange(value.toUpperCase());
    }
  };
</script>

<div>
  {#each inputs as input, index (index)}
    <input
      {disabled}
      bind:this={input}
      {placeholder}
      type="text"
      on:input={(e) => handleInput(e, index)}
      on:keydown={(e) => handleKeyDown(e, index)}
      on:keyup={() => handleKeyUp()}
      on:paste|preventDefault={handlePaste}
      maxlength="1"
      value={defaultValue?.[index] ?? ""}
    />
  {/each}
</div>

<style>
  div {
    display: flex;
    gap: var(--sl-spacing-x-small);
  }

  input {
    width: var(--sl-input-height-small);
    height: var(--sl-input-height-small);
    line-height: var(--sl-input-height-small);
    text-align: center;
    text-transform: uppercase;

    background-color: var(--sl-input-background-color);
    border-style: solid;
    border-color: var(--sl-input-border-color);
    border-width: var(--sl-input-border-width);
    border-radius: var(--sl-input-border-radius-small);
    font-family: var(--sl-input-font-family);
    font-weight: var(--sl-input-font-weight);
    font-size: var(--sl-input-font-size-small);
    color: var(--sl-input-color);
    outline: none;

    &:focus {
      background-color: var(--sl-input-background-color-focus);
      border-color: var(--sl-input-border-color-focus);
      outline-color: var(--sl-input-focus-ring-color);
      outline-offset: var(--sl-input-focus-ring-offset);
      color: var(--sl-input-color-focus);
    }

    &:hover {
      background-color: var(--sl-input-background-color-hover);
      border-color: var(--sl-input-border-color-hover);
      color: var(--sl-input-color-hover);
    }
  }

  input::placeholder {
    color: var(--sl-input-placeholder-color);
  }

  input::placeholder:disabled {
    color: var(--sl-input-placeholder-color-disabled);
  }

  input:disabled {
    background-color: var(--sl-input-background-color-disabled);
    border-color: var(--sl-input-border-color-disabled);
    color: var(--sl-input-color-disabled);
    cursor: not-allowed;
  }
</style>
