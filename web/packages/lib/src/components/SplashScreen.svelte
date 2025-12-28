<script lang="ts">
  import { onMount } from "svelte";
  import FullLogo from "./FullLogo.svelte";

  let { onComplete }: { onComplete?: () => void } = $props();

  let visible = $state(true);
  let startTime: number;
  let completed = false;

  onMount(() => {
    startTime = Date.now();

    const fallbackTimeout = setTimeout(() => {
      if (!completed) {
        handleCompletion();
      }
    }, 2500);

    return () => clearTimeout(fallbackTimeout);
  });

  const handleCompletion = () => {
    if (completed) return;
    completed = true;

    const elapsed = Date.now() - startTime;
    const remaining = Math.max(0, 2000 - elapsed);

    setTimeout(() => {
      visible = false;
      onComplete?.();
    }, remaining);
  };

  const handleAnimationEnd = () => {
    handleCompletion();
  };
</script>

{#if visible}
  <div class="splash-screen" onanimationend={handleAnimationEnd}>
    <div class="logo">
      <FullLogo />
    </div>
  </div>
{/if}

<style>
  .splash-screen {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: var(--wa-color-brand-50);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    animation: slide-in 0.5s ease-out;
  }

  @keyframes slide-in {
    from {
      transform: translateY(-100%);
    }
    to {
      transform: translateY(0);
    }
  }

  .logo {
    width: 50%;
    max-width: 600px;
    color: white;
  }
</style>
