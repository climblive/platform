<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import ProblemForm, { formSchema } from "@/forms/ProblemForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/button-group/button-group.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
  import type WaButton from "@awesome.me/webawesome/dist/components/button/button.js";
  import type { ProblemTemplate } from "@climblive/lib/models";
  import {
    createProblemMutation,
    getProblemsQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const problemsQuery = $derived(getProblemsQuery(contestId));

  let highestProblemNumber = $derived.by(() => {
    if (problemsQuery.data === undefined) {
      return undefined;
    } else if (problemsQuery.data.length > 0) {
      return Math.max(
        ...(problemsQuery.data?.map(({ number }) => number) ?? []),
      );
    } else {
      return 0;
    }
  });

  const createProblem = $derived(createProblemMutation(contestId));

  let addAnother = $state(false);
  let formKey = $state(0);
  let createBtn = $state<WaButton>();

  const handleSubmit = async (tmpl: Omit<ProblemTemplate, "pointsZone">) => {
    const shouldAddAnother = addAnother;
    addAnother = false;

    createProblem.mutate(
      { ...tmpl },
      {
        onSuccess: () => {
          if (shouldAddAnother) {
            formKey++;
          } else {
            navigate(`/admin/contests/${contestId}#problems`);
          }
        },
        onError: () => toastError("Failed to create problem."),
      },
    );
  };

  const handleDropdownSelect = (e: Event) => {
    const { item } = (e as CustomEvent<{ item: { value: string } }>).detail;
    if (item.value === "add-another") {
      addAnother = true;
      createBtn?.closest("form")?.requestSubmit();
    }
  };
</script>

{#if highestProblemNumber === undefined}
  <Loader />
{:else}
  {#key formKey}
    <ProblemForm
      submit={handleSubmit}
      data={{
        number: highestProblemNumber + 1,
        holdColorPrimary: "#000000",
        pointsTop: 100,
        flashBonus: 0,
      }}
      schema={formSchema}
    >
      <div class="controls">
        <wa-button
          size="small"
          type="button"
          appearance="plain"
          onclick={() => navigate(`/admin/contests/${contestId}#problems`)}
          >Cancel</wa-button
        >
        <wa-button-group>
          <wa-button
            bind:this={createBtn}
            size="small"
            type="submit"
            loading={createProblem.isPending}
            variant="neutral"
            >Create</wa-button
          >
          <wa-dropdown
            placement="bottom-end"
            onwa-select={handleDropdownSelect}
          >
            <wa-button
              slot="trigger"
              size="small"
              type="button"
              variant="neutral"
              loading={createProblem.isPending}
              with-caret
            ></wa-button>
            <wa-dropdown-item value="add-another"
              >Create and add another</wa-dropdown-item
            >
          </wa-dropdown>
        </wa-button-group>
      </div>
    </ProblemForm>
  {/key}
{/if}

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    justify-content: end;
  }
</style>
