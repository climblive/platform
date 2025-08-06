<script lang="ts">
  import { serialize } from "@awesome.me/webawesome";
  import { type Snippet } from "svelte";
  import * as z from "zod";

  type T = $$Generic<unknown>;

  interface Props {
    schema: z.ZodType<T, z.ZodTypeDef, T>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { schema, submit, children }: Props = $props();

  let form: HTMLFormElement | undefined = $state();

  const stash = new Map<string, string | null>();

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();

    if (!form) {
      return;
    }

    const data = serialize(form);
    const result = schema.safeParse(data);

    if (result.success) {
      submit(result.data);
    } else {
      for (const issue of result.error.issues) {
        setCustomValidity(issue.path, issue.message);
      }
    }

    form?.reportValidity();
  };

  const setCustomValidity = (path: (string | number)[], message: string) => {
    const name = `${path[0]}`;

    const input = form?.querySelector(`[name="${name}"]`) as
      | HTMLInputElement
      | null
      | undefined;

    if (!input) {
      return;
    }

    input.setCustomValidity(message);

    stash.set(input.name, input.getAttribute("hint"));

    input.setAttribute("hint", message);
  };

  const resetCustomValidation = () => {
    const inputs = form?.querySelectorAll(`[name]`) as
      | NodeListOf<HTMLInputElement>
      | undefined;

    if (!inputs) {
      return;
    }

    for (const input of inputs) {
      input.setCustomValidity("");

      const stashedHelpText = stash.get(input.name);

      switch (stashedHelpText) {
        case undefined:
          break;
        case null:
          input.removeAttribute("hint");
          break;
        default:
          input.setAttribute("hint", stashedHelpText);
      }

      stash.delete(input.name);
    }
  };
</script>

<form bind:this={form} onsubmit={handleSubmit} oninput={resetCustomValidation}>
  {@render children?.()}
</form>
