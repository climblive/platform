<script lang="ts" generics="T">
  import { serialize } from "@shoelace-style/shoelace";
  import { type Snippet } from "svelte";
  import * as z from "zod";

  interface Props {
    schema: z.ZodType<T, z.ZodTypeDef, T>;
    submit: (value: T) => void;
    children?: Snippet;
  }

  let { schema, submit, children }: Props = $props();

  let form: HTMLFormElement | undefined = $state();

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
    const input = form?.querySelector(`[name="${path}"]`) as
      | HTMLInputElement
      | null
      | undefined;
    input?.setCustomValidity(message);
  };

  const resetCustomValidation = () => {
    const inputs = form?.querySelectorAll(`[name]`) as
      | NodeListOf<HTMLInputElement>
      | undefined;

    if (!inputs) {
      return;
    }

    for (const input of inputs) {
      input?.setCustomValidity("");
    }
  };
</script>

<form bind:this={form} onsubmit={handleSubmit} oninput={resetCustomValidation}>
  {@render children?.()}
</form>
