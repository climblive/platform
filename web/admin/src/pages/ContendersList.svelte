<script lang="ts">
  import { GenericForm, name } from "@climblive/lib/forms";
  import {
    createContendersMutation,
    getContendersByContestQuery,
  } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/range/range.js";
  import type { CreateContendersArguments } from "node_modules/@climblive/lib/src/models/rest";
  import * as z from "zod";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contendersQuery = getContendersByContestQuery(contestId);
  const createContenders = createContendersMutation(contestId);

  let contenders = $derived($contendersQuery.data);

  let remainingCodes = $derived(
    contenders === undefined ? undefined : 500 - contenders.length,
  );

  const schema: z.ZodType<CreateContendersArguments> = $derived(
    z.object({
      number: z.coerce
        .number()
        .min(1)
        .max(remainingCodes ?? 0),
    }),
  );
</script>

{#if contenders}
  <GenericForm {schema} submit={(v) => $createContenders.mutate(v)}>
    <sl-range
      label="Number of registration codes"
      help-text={`You have ${remainingCodes} codes remaining`}
      min={0}
      max={remainingCodes}
      use:name={"number"}
      disabled={remainingCodes === undefined}
      step="10"
    ></sl-range>
    <sl-button
      size="small"
      type="submit"
      variant="primary"
      loading={$createContenders.isPending}
      >Create
    </sl-button>
  </GenericForm>

  <ul>
    {#each contenders as contender (contender.id)}
      <li>{contender.registrationCode}</li>
    {/each}
  </ul>
{/if}

<style>
  sl-range {
    --track-color-active: var(--sl-color-primary-600);
  }
</style>
