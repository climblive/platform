<script lang="ts">
  import imgUrl from "@/static/logo.svg";
</script>

<p>
  <img src={imgUrl} alt="ClimbLive" />
  <em>Brought to you by ClimbLive</em>
</p>

<style>
  p {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--sl-font-size-medium);
  }

  img {
    width: 1.25rem;
    width: 1.25rem;
    margin-inline: var(--sl-spacing-x-small);
  }
</style>
