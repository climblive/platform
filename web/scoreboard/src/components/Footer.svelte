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
    background-color: var(--sl-color-primary-300);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--sl-color-primary-950);
    font-size: var(--sl-font-size-small);
  }

  img {
    width: 1.5rem;
    width: 1.5rem;
    margin-right: var(--sl-spacing-x-small);
  }
</style>
