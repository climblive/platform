<script lang="ts">
  import "@awesome.me/webawesome/dist/components/spinner/spinner.js";
</script>

<main>
  <wa-spinner></wa-spinner>
</main>

<style>
  main {
    height: 100vh;
    display: flex;
    align-items: center;
  }

  wa-spinner {
    --track-width: 0.75rem;

    display: flex;
    font-size: 5rem;
    margin: 0 auto;
  }
</style>
