<script lang="ts">
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
</script>

<div role="status">
  <wa-icon name="check"></wa-icon> Saved
</div>

<style>
  div {
    display: flex;
    align-items: center;
    gap: var(--wa-space-2xs);
    font-size: var(--wa-font-size-s);
    color: var(--wa-color-success-fill-loud);
    font-weight: var(--wa-font-weight-bold);
    animation: hide 0s ease 2s forwards;
  }

  @keyframes hide {
    to {
      visibility: hidden;
    }
  }
</style>
