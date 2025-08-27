<script lang="ts">
  import "@awesome.me/webawesome/dist/styles/native.css";
</script>

<main>
  <section>
    <h1>Sorry!</h1>
    <p>
      Your browser version is outdated and may not support this application. We
      recommend you to upgrade your browser or borrow your friends phone.
    </p>
    <p>
      If you are using an iPhone or iPad, please ensure you <a
        href="https://support.apple.com/en-us/118575"
        >update to the latest version of iOS</a
      >. Please note that devices older than the iPhone XR may not be
      upgradable.
    </p>
    <button
      onclick={() => {
        const currentUrl = new URL(location.href);
        currentUrl.searchParams.set("force-compat", "true");

        location.href = currentUrl.toString();
      }}>Try anyway</button
    >
  </section>
</main>

<style>
  main {
    padding: var(--wa-space-m);
  }

  section {
    padding: var(--wa-space-m);
    background-color: var(--wa-surface-raised);
  }
</style>
