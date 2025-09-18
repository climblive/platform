<script lang="ts" module>
  import * as z from "zod";

  export const formSchema = z.object({
    location: z.string().optional(),
    seriesId: z.coerce.number().optional(),
    name: z.string().min(1),
    description: z.string().optional(),
    qualifyingProblems: z.coerce.number().min(0).max(65536),
    finalists: z.coerce.number().min(0).max(65536),
    rules: z.string().optional(),
    gracePeriod: z.coerce.number().min(0).max(60),
  });

  export const minuteInNanoseconds = 60 * 1_000_000_000;
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button-group/button-group.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/textarea/textarea.js";
  import type WaTextarea from "@awesome.me/webawesome/dist/components/textarea/textarea.js";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { Contest } from "@climblive/lib/models";
  import { Editor } from "@tiptap/core";
  import { StarterKit } from "@tiptap/starter-kit";
  import { onDestroy, onMount, type Snippet } from "svelte";

  type T = $$Generic<Partial<Contest>>;

  interface Props {
    data: Partial<T>;
    schema: z.ZodType<T, z.ZodTypeDef, T>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { data, schema, submit, children }: Props = $props();

  let element: HTMLElement | undefined = $state();
  let rules: WaTextarea | undefined = $state();
  let editor = $state<Editor>();

  onMount(() => {
    editor = new Editor({
      element: element,
      extensions: [StarterKit],
      content: data.rules,

      onUpdate: ({ editor }) => {
        if (rules) {
          rules.value = editor.getHTML();
        }
      },
    });
  });
  onDestroy(() => {
    editor?.destroy();
  });
</script>

<GenericForm {schema} {submit}>
  <fieldset>
    <wa-input
      size="small"
      {@attach name("name")}
      label="Name"
      type="text"
      required
      value={data.name}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("description")}
      label="Description"
      type="text"
      value={data.description}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("location")}
      label="Location"
      type="text"
      value={data.location}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("finalists")}
      label="Finalists"
      hint="Number of contenders that will proceed to the finals."
      type="number"
      required
      value={data.finalists}
      min={0}
      max={65536}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("qualifyingProblems")}
      label="Number of qualifying problems"
      hint="Number of the hardest problems that will count towards each contender's score."
      type="number"
      required
      value={data.qualifyingProblems}
      min={0}
      max={65536}
    ></wa-input>
    <wa-input
      size="small"
      {@attach name("gracePeriod")}
      label="Grace period (minutes)"
      hint="Extra time after the end of the contest during which contenders can enter their last results."
      type="number"
      required
      min={0}
      max={60}
      value={Math.floor((data.gracePeriod ?? 0) / minuteInNanoseconds)}
    >
    </wa-input>
    <div>
      <wa-textarea
        size="small"
        {@attach name("rules")}
        label="Rules"
        value={data.rules ?? ""}
        bind:this={rules}
      ></wa-textarea>
      {#if editor}
        <wa-button-group>
          <wa-button
            size="small"
            onclick={(e: MouseEvent) => {
              e.preventDefault();
              editor?.chain().focus().toggleHeading({ level: 1 }).run();
            }}
            class:active={editor.isActive("heading", { level: 1 })}
          >
            <wa-icon name="heading"></wa-icon>
          </wa-button>
          <wa-button
            size="small"
            onclick={(e: MouseEvent) => {
              e.preventDefault();
              editor?.chain().focus().setParagraph().run();
            }}
            class:active={editor.isActive("paragraph")}
          >
            <wa-icon name="paragraph"></wa-icon>
          </wa-button>
          <wa-button
            size="small"
            onclick={(e: MouseEvent) => {
              e.preventDefault();
              editor?.chain().focus().setItalic().run();
            }}
            class:active={editor.isActive("paragraph")}
          >
            <wa-icon name="italic"></wa-icon>
          </wa-button>
          <wa-button
            size="small"
            onclick={(e: MouseEvent) => {
              e.preventDefault();
              editor?.chain().focus().setBold().run();
            }}
            class:active={editor.isActive("paragraph")}
          >
            <wa-icon name="bold"></wa-icon>
          </wa-button>
          <wa-button
            size="small"
            onclick={(e: MouseEvent) => {
              e.preventDefault();
              editor?.chain().focus().setUnderline().run();
            }}
            class:active={editor.isActive("paragraph")}
          >
            <wa-icon name="underline"></wa-icon>
          </wa-button>
          <wa-button
            size="small"
            onclick={(e: MouseEvent) => {
              e.preventDefault();
              editor?.chain().focus().setStrike().run();
            }}
            class:active={editor.isActive("paragraph")}
          >
            <wa-icon name="strikethrough"></wa-icon>
          </wa-button>
        </wa-button-group>
      {/if}
      <div bind:this={element} class="rules"></div>
    </div>
    {@render children?.()}
  </fieldset>
</GenericForm>

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

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
