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
    } catch {}

    if (shouldSkipSplash) {
      onComplete();
      showSplash = false;
      return;
    }

    const fallbackTimeout = setTimeout(() => {
      try {
        sessionStorage.setItem("splashShown", "true");
      } catch {}

      onComplete();
    }, 2_000);

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
    <wa-spinner size="xl"></wa-spinner>
  </div>
{/if}

<style>
  .splash-screen,
  .spinner-screen {
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .splash-screen {
    background-color: var(--wa-color-brand-fill-loud);
  }

  .logo {
    width: 50%;
    color: white;
    animation: slide-in 0.5s var(--wa-transition-easing);
  }

  wa-spinner {
    font-size: 5rem;
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
