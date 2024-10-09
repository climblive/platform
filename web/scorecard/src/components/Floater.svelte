<script lang="ts">
  export let order: number;
</script>

<div
  style="transform: translateY(calc({order} * 2.25rem + {order} * var(--sl-spacing-x-small)))"
>
  <slot />
</div>

<style>
  div {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    transition: transform 0.5s ease-out;
  }
</style>
