<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button-group/button-group.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/textarea/textarea.js";
  import type WaTextarea from "@awesome.me/webawesome/dist/components/textarea/textarea.js";
  import { name } from "@climblive/lib/forms";
  import { type ChainedCommands, Editor } from "@tiptap/core";
  import { StarterKit } from "@tiptap/starter-kit";
  import { onDestroy, onMount } from "svelte";

  interface Props {
    rules: string | undefined;
  }

  let { rules }: Props = $props();

  let element: HTMLElement | undefined = $state();
  let rulesElement: WaTextarea | undefined = $state();
  let editor = $state<Editor>();

  onMount(() => {
    editor = new Editor({
      element: element,
      extensions: [StarterKit],
      content: rules,

      onUpdate: ({ editor }) => {
        if (rulesElement) {
          rulesElement.value = editor.getHTML();
        }
      },
    });
  });
  onDestroy(() => {
    editor?.destroy();
  });
</script>

{#snippet richTextModifier(
  editor: Editor,
  action: (chain: ChainedCommands) => ChainedCommands,
  isActive: (editor: Editor) => boolean,
  iconName: string,
)}
  <wa-button
    size="small"
    onclick={(e: MouseEvent) => {
      e.preventDefault();

      const chain = editor.chain().focus();
      action(chain).run();
    }}
    class:active={isActive(editor)}
  >
    <wa-icon name={iconName}></wa-icon>
  </wa-button>
{/snippet}

<div>
  <wa-textarea
    size="small"
    {@attach name("rules")}
    label="Rules"
    value={rules ?? ""}
    bind:this={rulesElement}
  ></wa-textarea>
  {#if editor}
    <wa-button-group>
      {@render richTextModifier(
        editor,
        (chain) => chain.toggleHeading({ level: 2 }),
        (editor) => editor.isActive("heading", { level: 2 }),
        "heading",
      )}
      {@render richTextModifier(
        editor,
        (chain) => chain.setParagraph(),
        (editor) => editor.isActive("paragraph"),
        "paragraph",
      )}
      {@render richTextModifier(
        editor,
        (chain) => chain.toggleItalic(),
        (editor) => editor.isActive("italic"),
        "italic",
      )}
      {@render richTextModifier(
        editor,
        (chain) => chain.toggleBold(),
        (editor) => editor.isActive("bold"),
        "bold",
      )}
      {@render richTextModifier(
        editor,
        (chain) => chain.toggleUnderline(),
        (editor) => editor.isActive("underline"),
        "underline",
      )}
      {@render richTextModifier(
        editor,
        (chain) => chain.toggleStrike(),
        (editor) => editor.isActive("strike"),
        "strikethrough",
      )}
      {@render richTextModifier(
        editor,
        (chain) => chain.toggleBulletList(),
        (editor) => editor.isActive("bulletList"),
        "list-ul",
      )}
      {@render richTextModifier(
        editor,
        (chain) => chain.toggleOrderedList(),
        (editor) => editor.isActive("orderedList"),
        "list-ol",
      )}
    </wa-button-group>
  {/if}
  <div bind:this={element} class="rules"></div>
</div>

<style>
  .rules :global(.tiptap) {
    min-height: 10rem;
    border: var(--wa-form-control-border-style)
      var(--wa-form-control-border-width) var(--wa-form-control-border-color);
    border-radius: var(--wa-form-control-border-radius);
    padding: 0 var(--wa-form-control-padding-inline);
  }

  wa-textarea::part(base) {
    display: none;
  }

  wa-button-group {
    margin-block-end: var(--wa-space-xs);
  }
</style>
