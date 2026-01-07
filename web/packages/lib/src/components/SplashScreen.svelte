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
      showSplash = false;
      return;
    }

    const fallbackTimeout = setTimeout(() => {
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
    <wa-spinner size="xl"></wa-spinner>
  </div>
{/if}

<style>
  .splash-screen,
  .spinner-screen {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .splash-screen {
    background-color: var(--wa-color-brand-fill-loud);
  }

  .splash-screen {
    padding-bottom: 10vh;
  }

  .logo {
    width: 50%;
    max-width: 600px;
    color: white;
    animation: slide-in 0.5s ease-out;
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
