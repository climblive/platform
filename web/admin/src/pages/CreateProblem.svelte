<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import ProblemForm, { formSchema } from "@/forms/ProblemForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
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
      {#snippet children(form)}
        <div class="controls">
          <wa-button
            size="small"
            type="button"
            appearance="plain"
            onclick={() => navigate(`/admin/contests/${contestId}#problems`)}
            >Cancel</wa-button
          >
          <div class="split-button">
            <wa-button
              class="create-btn"
              size="small"
              type="submit"
              loading={createProblem.isPending}
              variant="neutral"
              onclick={() => (addAnother = false)}
              >Create
            </wa-button>
            <wa-dropdown placement="bottom-end">
              <wa-button
                class="caret-btn"
                slot="trigger"
                size="small"
                type="button"
                variant="neutral"
                loading={createProblem.isPending}
                with-caret
              ></wa-button>
              <wa-dropdown-item
                onclick={() => {
                  addAnother = true;
                  form?.requestSubmit();
                }}>Create and add another</wa-dropdown-item
              >
            </wa-dropdown>
          </div>
        </div>
      {/snippet}
    </ProblemForm>
  {/key}
{/if}

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    justify-content: end;
  }

  .split-button {
    display: inline-flex;
  }

  .split-button .create-btn::part(base) {
    border-start-end-radius: 0;
    border-end-end-radius: 0;
  }

  .split-button .caret-btn::part(base) {
    border-start-start-radius: 0;
    border-end-start-radius: 0;
  }
</style>
