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

<a href={code ? `/failsafe/${code}` : "/failsafe"}
  ><button>Use simplified version</button></a
>
