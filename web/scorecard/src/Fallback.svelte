<script lang="ts">
  import "@awesome.me/webawesome/dist/styles/native.css";
</script>

<main>
  <section>
    <p>Your browser is too old.</p>
    <button>Failsafe</button>
  </section>
</main>

<style>
  main {
    width: 100%;
  }

  section {
    margin: var(--wa-space-m);
    padding: var(--wa-space-m);
    background-color: var(--wa-surface-raised);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-neutral-border-normal);
    border-radius: var(--wa-border-radius-m);
  }
</style>
