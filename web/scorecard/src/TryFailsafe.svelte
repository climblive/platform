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
  Alternatively you can use a simplified version of the app that includes only
  the most basic features.
  <br />

  <a href={code ? `/failsafe/${code}` : "/failsafe"}
    ><button>Use simplified version</button></a
  >
</p>

<style>
  p {
    margin-block-start: var(--wa-space-m);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-success-border-loud);
    padding: var(--wa-space-m);
    border-radius: var(--wa-border-radius-m);
  }

  a {
    display: block;
    margin-block-start: var(--wa-space-m);
  }
</style>
