<script lang="ts">
  import { serialize } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import SlInput from "@shoelace-style/shoelace/dist/components/input/input.js";
  import "@shoelace-style/shoelace/dist/components/option/option.js";
  import "@shoelace-style/shoelace/dist/components/select/select.js";
  import type SlSelect from "@shoelace-style/shoelace/dist/components/select/select.js";
  import { createEventDispatcher, getContext } from "svelte";
  import type { Readable } from "svelte/store";
  import {
    registrationFormSchema,
    type RegistrationFormData,
  } from "@climblive/shared/models";
  import { getCompClassesQuery } from "@climblive/shared/queries";
  import type { ScorecardSession } from "@/types";

  const dispatch = createEventDispatcher<{ submit: RegistrationFormData }>();

  export let data: Partial<RegistrationFormData>;

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  $: compClassesQuery = getCompClassesQuery($session.contestId);

  let form: HTMLFormElement;
  const controls: {
    name?: SlInput;
    club?: SlInput;
    compClassId?: SlSelect;
  } = {};

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
    on:submit|preventDefault={handleSubmit}
    on:sl-input={resetCustomValidation}
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
      bind:this={controls.club}
      size="small"
      name="club"
      label="Club name"
      type="text"
      use:value={data.club}
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
        <sl-option value={compClass.id}>{compClass.name}</sl-option>
      {/each}
    </sl-select>
    <slot />
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
