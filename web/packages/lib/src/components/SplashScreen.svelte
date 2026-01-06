<script lang="ts">
  import { onMount } from "svelte";
  import SplashLogo from "./SplashLogo.svelte";

  let { onComplete }: { onComplete?: () => void } = $props();

  let visible = $state(true);
  let startTime: number;
  let completed = false;

  onMount(() => {
    try {
      const hasShown = sessionStorage.getItem("splashShown");
      
      if (hasShown) {
        visible = false;
        onComplete?.();
        return;
      }
    } catch {
      // sessionStorage may be unavailable in private browsing or when disabled
    }

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
      try {
        sessionStorage.setItem("splashShown", "true");
      } catch {
        // sessionStorage may be unavailable in private browsing or when disabled
      }
      visible = false;
      onComplete?.();
    }, remaining);
  };

  const handleAnimationEnd = () => {
    handleCompletion();
  };
</script>

{#if visible}
  <div class="splash-screen">
    <div class="logo" onanimationend={handleAnimationEnd}>
      <SplashLogo />
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
  }

  .logo {
    width: 50%;
    max-width: 600px;
    color: white;
    animation: slide-in 0.5s ease-out;
  }

  @keyframes slide-in {
    from {
      transform: translateX(-100%);
    }
    to {
      transform: translateX(0);
    }
  }
</style>
