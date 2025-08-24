<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { FullLogo } from "@climblive/lib/components";
  import { getContext, onMount } from "svelte";
  import type { Authenticator } from "./authenticator.svelte";

  let print = $state(false);

  const authenticator = getContext<Authenticator>("authenticator");

  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get("print") !== null) {
      print = true;
    }
  });
</script>

{#if !print}
  <header>
    <div>
      <p class="logo">
        <FullLogo />
      </p>
      <wa-button
        size="small"
        appearance="outlined"
        onclick={authenticator.logout}
      >
        Sign out<wa-icon slot="start" name="right-from-bracket"
        ></wa-icon></wa-button
      >
    </div>
  </header>
{/if}

<style>
  header {
    background-color: var(--wa-color-surface-lowered);
  }

  div {
    margin: 0 auto;
    max-width: 1024px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-inline-end: var(--wa-space-m);
    height: 3.5rem;
    gap: var(--wa-space-xl);
  }

  .logo {
    text-align: left;
    height: var(--wa-font-size-xl);
    color: var(--wa-color-text-normal);
    padding-left: var(--wa-space-xs);
    flex-shrink: 0;
    margin-inline-start: var(--wa-space-xs);
  }

  @media print {
    header {
      display: none;
    }
  }
</style>
