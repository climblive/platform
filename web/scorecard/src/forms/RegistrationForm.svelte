<script lang="ts">
  import { preventDefault } from "svelte/legacy";

  import type { ScorecardSession } from "@/types";
  import {
    registrationFormSchema,
    type RegistrationFormData,
  } from "@climblive/lib/models";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import { serialize } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import SlInput from "@shoelace-style/shoelace/dist/components/input/input.js";
  import "@shoelace-style/shoelace/dist/components/option/option.js";
  import "@shoelace-style/shoelace/dist/components/select/select.js";
  import type SlSelect from "@shoelace-style/shoelace/dist/components/select/select.js";
  import { isAfter } from "date-fns";
  import { createEventDispatcher, getContext, type Snippet } from "svelte";
  import type { Readable } from "svelte/store";

  const dispatch = createEventDispatcher<{ submit: RegistrationFormData }>();

  interface Props {
    data: Partial<RegistrationFormData>;
    children?: Snippet;
  }

  let { data, children }: Props = $props();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  let compClassesQuery = $derived(getCompClassesQuery($session.contestId));

  let form: HTMLFormElement = $state();
  const controls: {
    name?: SlInput;
    clubName?: SlInput;
    compClassId?: SlSelect;
  } = $state({});

  const handleSubmit = () => {
    const data = serialize(form);
    const result = registrationFormSchema.safeParse(data);

    if (result.success) {
      dispatch("submit", result.data);
    } else {
      for (const issue of result.error.issues) {
        setCustomValidity(issue.path, issue.message);
      }
    }

    form.reportValidity();
  };

  const setCustomValidity = (path: (string | number)[], message: string) => {
    for (const [key, input] of Object.entries(controls)) {
      if (key === path[0]) {
        input?.setCustomValidity(message);
        return;
      }
    }
  };

  const resetCustomValidation = () => {
    for (const input of Object.values(controls)) {
      input?.setCustomValidity("");
    }
  };

  const value = (node: HTMLElement, value: string | number | undefined) => {
    node.setAttribute("value", value?.toString() ?? "");

    return {
      update(value: string | number | undefined) {
        node.setAttribute("value", value?.toString() ?? "");
      },
    };
  };
</script>

{#if $compClassesQuery.data}
  <form
    bind:this={form}
    onsubmit={preventDefault(handleSubmit)}
    onsl-input={resetCustomValidation}
  >
    <sl-input
      bind:this={controls.name}
      size="small"
      name="name"
      label="Full name"
      type="text"
      required
      use:value={data.name}
    ></sl-input>
    <sl-input
      bind:this={controls.clubName}
      size="small"
      name="clubName"
      label="Club name"
      type="text"
      use:value={data.clubName}
    ></sl-input>
    <sl-select
      bind:this={controls.compClassId}
      size="small"
      name="compClassId"
      label="Competition class"
      required
      use:value={data.compClassId}
    >
      {#each $compClassesQuery.data as compClass}
        <sl-option
          value={compClass.id}
          disabled={isAfter(new Date(), compClass.timeEnd)}
          >{compClass.name}</sl-option
        >
      {/each}
    </sl-select>
    {@render children?.()}
  </form>
{/if}

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-small);
    padding: var(--sl-spacing-medium);
  }
</style>
