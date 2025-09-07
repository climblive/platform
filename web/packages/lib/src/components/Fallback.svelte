<script lang="ts">
  import { type Component } from "svelte";
  import { importNativeStyles } from "./styles";

  type Props = {
    missingFeatures?: string[];
    alternative?: Component;
  };

  const { missingFeatures = [], alternative: Alternative }: Props = $props();

  let tapCount = $state(0);
  let showMissingFeatures = $derived(tapCount >= 5);

  const handleTap = () => {
    tapCount += 1;
  };

  const handleForce = () => {
    const url = new URL(window.location.href);
    url.searchParams.set("compat", "ignore");

    window.location.href = url.toString();
  };
</script>

{#await importNativeStyles() then styles}
  <svelte:component this={styles.default} />
{/await}
<main>
  <section>
    <h1 onclickcapture={handleTap}>Sorry!</h1>
    <p>
      Your browser version is outdated and will most likely not support this
      app. We recommend you to upgrade<sup>*</sup> your browser or borrow your friends
      phone.
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

    {#if Alternative}
      <Alternative />
    {/if}

    <p>
      <small>
        You may try your luck and
        <a href="#" onclick={handleForce} class="wa-danger wa-size-s"
          >continue anyway</a
        >, but be aware that the app might not work as expected.
      </small>
    </p>

    <p>
      <small>
        <sup>*</sup>If you are using an iPhone or iPad, please ensure you
        <a href="https://support.apple.com/en-us/118575"
          >update to the latest version of iOS</a
        >. Please note that devices older than the iPhone XR may not be
        upgradable.
      </small>
    </p>
  </section>
</main>

<style>
  h1 {
    user-select: none;
  }

  main {
    padding: var(--wa-space-m);
  }

  section {
    padding: var(--wa-space-m);
  }
</style>
