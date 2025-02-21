<script lang="ts">
  import { createContestMutation } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/input/input.js";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const createContest = createContestMutation(organizerId);
</script>

<form>
  <sl-input label="Name"></sl-input>
</form>
