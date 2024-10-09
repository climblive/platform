<script lang="ts">
  import "@shoelace-style/shoelace/dist/components/spinner/spinner.js";
</script>

<main>
  <sl-spinner></sl-spinner>
</main>

<style>
  main {
    height: 100vh;
    display: flex;
    align-items: center;
  }

  sl-spinner {
    --track-width: 0.75rem;

    display: flex;
    font-size: 5rem;
    margin: 0 auto;
  }
</style>
