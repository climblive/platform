<script lang="ts">
  import imgUrl from "@/static/logo.svg";
</script>

<footer>
  <img src={imgUrl} alt="ClimbLive" />
  Powered by ClimbLive
</footer>

<style>
  footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 3rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--sl-color-neutral-50);
    font-size: var(--sl-font-size-medium);
  }

  img {
    width: 1.25rem;
    width: 1.25rem;
    margin-right: var(--sl-spacing-x-small);
  }
</style>
