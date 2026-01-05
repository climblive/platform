<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button-group/button-group.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { name } from "@climblive/lib/forms";
  import { type ChainedCommands, Editor } from "@tiptap/core";
  import { StarterKit } from "@tiptap/starter-kit";
  import { onDestroy, onMount } from "svelte";

  interface Props {
    info: string | undefined;
  }

  const { info }: Props = $props();

  let infoElement: HTMLElement | undefined = $state();
  let hiddenElement: HTMLInputElement | undefined = $state();
  let editorState = $state<{ editor: Editor | null }>({ editor: null });

  const id = $props.id();

  onMount(() => {
    editorState.editor = new Editor({
      element: infoElement,
      extensions: [StarterKit],
      content: info,
      editorProps: {
        attributes: {
          "aria-labelledby": id,
          "aria-multiline": "true",
        },
      },
      onUpdate: ({ editor }) => {
        if (hiddenElement) {
          hiddenElement.value = editor.getHTML();
        }
      },
      onTransaction: ({ editor }) => {
        editorState = { editor };
      },
    });
  });

  onDestroy(() => {
    editorState.editor?.destroy();
  });
</script>

{#snippet richTextModifier(options: {
  editor: Editor;
  action: (chain: ChainedCommands) => ChainedCommands;
  isActive: (editor: Editor) => boolean;
  iconName: string;
})}
  <wa-button
    size="small"
    onclick={(e: MouseEvent) => {
      e.preventDefault();

      const chain = options.editor.chain().focus();
      options.action(chain).run();
    }}
    appearance={options.isActive(options.editor)
      ? "filled-outlined"
      : "outlined"}
  >
    <wa-icon name={options.iconName}></wa-icon>
  </wa-button>
{/snippet}

<div>
  <!-- svelte-ignore a11y_label_has_associated_control -->
  <label {id}>General info</label>
  <input type="hidden" {@attach name("info")} bind:this={hiddenElement} />
  {#if editorState.editor}
    <wa-button-group>
      {@render richTextModifier({
        editor: editorState.editor,
        action: (chain: ChainedCommands) => chain.setHeading({ level: 2 }),
        isActive: (editor: Editor) => editor.isActive("heading", { level: 2 }),
        iconName: "heading",
      })}
      {@render richTextModifier({
        editor: editorState.editor,
        action: (chain: ChainedCommands) => chain.setParagraph(),
        isActive: (editor: Editor) => editor.isActive("paragraph"),
        iconName: "paragraph",
      })}
    </wa-button-group>

    <wa-button-group>
      {@render richTextModifier({
        editor: editorState.editor,
        action: (chain: ChainedCommands) => chain.toggleItalic(),
        isActive: (editor: Editor) => editor.isActive("italic"),
        iconName: "italic",
      })}
      {@render richTextModifier({
        editor: editorState.editor,
        action: (chain: ChainedCommands) => chain.toggleBold(),
        isActive: (editor: Editor) => editor.isActive("bold"),
        iconName: "bold",
      })}
      {@render richTextModifier({
        editor: editorState.editor,
        action: (chain: ChainedCommands) => chain.toggleUnderline(),
        isActive: (editor: Editor) => editor.isActive("underline"),
        iconName: "underline",
      })}
      {@render richTextModifier({
        editor: editorState.editor,
        action: (chain: ChainedCommands) => chain.toggleStrike(),
        isActive: (editor: Editor) => editor.isActive("strike"),
        iconName: "strikethrough",
      })}
    </wa-button-group>

    <wa-button-group>
      {@render richTextModifier({
        editor: editorState.editor,
        action: (chain: ChainedCommands) => chain.toggleBulletList(),
        isActive: (editor: Editor) => editor.isActive("bulletList"),
        iconName: "list-ul",
      })}
      {@render richTextModifier({
        editor: editorState.editor,
        action: (chain: ChainedCommands) => chain.toggleOrderedList(),
        isActive: (editor: Editor) => editor.isActive("orderedList"),
        iconName: "list-ol",
      })}
    </wa-button-group>
  {/if}

  <div bind:this={infoElement} class="info"></div>
</div>

<style>
  .info :global(.tiptap) {
    min-height: 10rem;
    border: var(--wa-form-control-border-style)
      var(--wa-form-control-border-width) var(--wa-form-control-border-color);
    border-radius: var(--wa-form-control-border-radius);
    padding: 0 var(--wa-form-control-padding-inline);
  }

  wa-button-group {
    margin-block-end: var(--wa-space-xs);

    &:not(:last-of-type) {
      margin-inline-end: var(--wa-space-xs);
    }

    &::part(base) {
      gap: 0.0625rem;
    }
  }

  wa-button::part(base) {
    border-color: var(--wa-color-border-loud);
  }

  label {
    display: block;
    font-size: var(--wa-font-size-s);
    color: var(--wa-form-control-label-color);
    font-weight: var(--wa-form-control-label-font-weight);
    line-height: var(--wa-form-control-label-line-height);
    margin-block-start: 0;
    margin-block-end: 0.5em;
  }
</style>
