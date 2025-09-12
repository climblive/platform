<script lang="ts">
  import { onMount } from "svelte";

  let code = $state<string>();

  const extractCodeFromPath = () => {
    const match = window.location.pathname.match(/\/([A-Z0-9]{8})/i);
    return match ? match[1] : null;
  };

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
