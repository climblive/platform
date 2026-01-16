<script lang="ts">
  import { onMount } from "svelte";
  import SplashLogo from "./SplashLogo.svelte";

  const { onComplete }: { onComplete: () => void } = $props();

  let showSplash = $state(true);

  onMount(() => {
    let shouldSkipSplash = false;

    try {
      const hasShown = sessionStorage.getItem("splashShown");

      if (hasShown === "true") {
        shouldSkipSplash = true;
      }
    } catch {
      /* sessionStorage may be unavailable (private browsing, disabled storage) */
    }

    if (shouldSkipSplash) {
      onComplete();
      showSplash = false;
      return;
    }

    const fallbackTimeout = setTimeout(() => {
      try {
        sessionStorage.setItem("splashShown", "true");
      } catch {
        /* sessionStorage may be unavailable (private browsing, disabled storage) */
      }

      onComplete();
    }, 1_500);

    return () => clearTimeout(fallbackTimeout);
  });
</script>

{#if showSplash}
  <div class="splash-screen">
    <div class="logo">
      <SplashLogo />
    </div>
  </div>
{:else}
  <div class="spinner-screen">
    <wa-spinner></wa-spinner>
  </div>
{/if}

<style>
  .splash-screen,
  .spinner-screen {
    width: 100%;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .splash-screen {
    background-color: var(--wa-color-brand-fill-loud);
  }

  .logo {
    width: min(20rem, 50%);
    color: white;
    animation: zoom-in 1s cubic-bezier(0.68, -0.55, 0.265, 1.55);
    margin-block-end: 10rem;
  }

  wa-spinner {
    font-size: 5rem;
  }

  @keyframes zoom-in {
    0% {
      transform: scale(0.3);
      opacity: 0;
    }
    50% {
      opacity: 1;
    }
    70% {
      transform: scale(1.05);
    }
    100% {
      transform: scale(1);
    }
  }
</style>
