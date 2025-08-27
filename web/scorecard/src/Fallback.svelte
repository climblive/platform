<script lang="ts">
  import "@awesome.me/webawesome/dist/styles/native.css";
</script>

<main>
  <section>
    <h1>Sorry!</h1>
    <p>
      Your browser is too old to run this app. Upgrade you browser to the latest
      version or borrow your friends phone.
    </p>
    <p>
      If you are using an iPhone or iPad you'll need to <a
        href="https://support.apple.com/en-us/118575"
        >update to the latest version of iOS</a
      >. Note that if your device is older than iPhone XR you wont be able to
      upgrade to the latest version.
    </p>
    <button>Try failsafe version</button>
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
