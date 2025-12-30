<script lang="ts">
  import { extractCodeFromPath } from "@climblive/lib/utils";
  import { onMount } from "svelte";

  let code = $state<string>();

  onMount(() => {
    const extractedCode = extractCodeFromPath();

    if (extractedCode) {
      code = extractedCode;
    }
  });
</script>

<p>
  If none of that is possible, you can use a basic version of the app instead.

  <a href={code ? `/failsafe/${code}` : "/failsafe"}
    ><button class="wa-success">Use basic version</button></a
  >
</p>

<style>
  a {
    display: block;
    margin-block-start: var(--wa-space-m);
  }

  button {
    width: 100%;
  }
</style>
