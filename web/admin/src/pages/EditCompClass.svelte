<script lang="ts">
  import CompClassForm from "@/forms/CompClassForm.svelte";
  import type { CompClassPatch } from "@climblive/lib/models";
  import {
    getCompClassQuery,
    patchCompClassMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import { navigate } from "svelte-routing";
  import * as z from "zod";

  const twelveHours = 12 * 60 * 60 * 1_000;

  interface Props {
    compClassId: number;
  }

  let { compClassId }: Props = $props();

  const formSchema: z.ZodType<CompClassPatch> = z
    .object({
      name: z.string().min(1),
      description: z.string().optional(),
      timeBegin: z.coerce.date(),
      timeEnd: z.coerce.date(),
    })
    .superRefine((data, ctx) => {
      if (data.timeEnd.getTime() - data.timeBegin.getTime() > twelveHours) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Total duration must not exceed 12 hours",
          path: ["timeEnd"],
        });
      }

      if (data.timeEnd <= data.timeBegin) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Time must follow chronological order",
          path: ["timeEnd"],
        });
      }
    });

  const compClassQuery = getCompClassQuery(compClassId);
  const patchCompClass = patchCompClassMutation(compClassId);

  const compClass = $derived($compClassQuery.data);

  const handleSubmit = async (patch: CompClassPatch) => {
    $patchCompClass.mutate(patch, {
      onSuccess: (compClass) =>
        navigate(`/admin/contests/${compClass.contestId}#comp-classes`),
      onError: () => toastError("Failed to save comp class."),
    });
  };
</script>

<CompClassForm
  submit={handleSubmit}
  data={{
    ...compClass,
  }}
  schema={formSchema}
>
  <div class="controls">
    <sl-button
      size="small"
      type="button"
      variant="text"
      onclick={history.back()}>Cancel</sl-button
    >
    <sl-button
      size="small"
      type="submit"
      loading={$patchCompClass.isPending}
      variant="primary"
      >Save
    </sl-button>
  </div>
</CompClassForm>
