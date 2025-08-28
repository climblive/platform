<script lang="ts">
  import type { Component } from "svelte";

  type Props = {
    missingFeatures?: string[];
    app: Component;
    styles: Component;
  };

  const { missingFeatures = [], app: App, styles: Styles }: Props = $props();

  let force = $state(false);
  let showMissingFeatures = $state(false);
  let tapCount = $state(0);

  const handleTap = () => {
    tapCount += 1;

    if (tapCount >= 5) {
      showMissingFeatures = true;
    }
  };
</script>

{#if force}
  <App></App>
{:else}
  <Styles />
  <main>
    <section>
      <h1 onclickcapture={handleTap}>Sorry!</h1>
      <p>
        Your browser version is outdated and may not support this application.
        We recommend you to upgrade your browser or borrow your friends phone.
      </p>

      {#if showMissingFeatures && missingFeatures.length > 0}
        <p>
          {#each missingFeatures as feature, index (index)}
            {#if index !== 0}
              ,&nbsp
            {/if}
            <code>{feature}</code>
          {/each}.
        </p>
      {/if}

      <button onclick={() => (force = true)} class="wa-danger wa-size-s"
        >Continue anyway</button
      >

      <p>
        If you are using an iPhone or iPad, please ensure you <a
          href="https://support.apple.com/en-us/118575"
          >update to the latest version of iOS</a
        >. Please note that devices older than the iPhone XR may not be
        upgradable.
      </p>
    </section>
  </main>
{/if}

<style>
  h1 {
    user-select: none;
  }

  main {
    padding: var(--wa-space-m);
  }

  button {
    margin-block-end: var(--wa-space-l);
  }

  section {
    padding: var(--wa-space-m);
    background-color: var(--wa-surface-raised);
  }
</style>
