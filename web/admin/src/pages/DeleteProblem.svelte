<script lang="ts">
  import { deleteProblemMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";

  type Props = {
    problemID: number;
    children: Snippet<[{ onDelete: () => void }]>;
  };

  let { problemID, children }: Props = $props();

  const deleteProblem = deleteProblemMutation(problemID);

  const handleDelete = async () => {
    console.log("Deleting problem with ID:", problemID);

    $deleteProblem.mutate(undefined, {
      onError: () => toastError("Failed to delete problem."),
    });
  };
</script>

{@render children({ onDelete: handleDelete })}
