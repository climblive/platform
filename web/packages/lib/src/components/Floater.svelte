<script lang="ts">
  export let order: number;
</script>

<div style="--order: {order}">
  <slot />
</div>

<style>
  div {
    position: absolute;
    left: 0;
    right: 0;
    transition: transform 0.2s linear;
    transform: translateY(
      calc(
        (var(--order) - var(--page-index) * var(--page-size)) *
          (2.25rem + var(--sl-spacing-x-small))
      )
    );
  }
</style>
