<div class="loader"></div>

<style>
  .loader {
    width: 60px;
    aspect-ratio: 4;
    --gradient: no-repeat
      radial-gradient(
        circle closest-side,
        var(--wa-color-brand-fill-loud) 90%,
        #0000
      );
    background:
      var(--gradient) 0% 50%,
      var(--gradient) 50% 50%,
      var(--gradient) 100% 50%;
    background-size: calc(100% / 3) 100%;
    animation: balls 1s infinite linear;
  }
  @keyframes balls {
    33% {
      background-size:
        calc(100% / 3) 0%,
        calc(100% / 3) 100%,
        calc(100% / 3) 100%;
    }
    50% {
      background-size:
        calc(100% / 3) 100%,
        calc(100% / 3) 0%,
        calc(100% / 3) 100%;
    }
    66% {
      background-size:
        calc(100% / 3) 100%,
        calc(100% / 3) 100%,
        calc(100% / 3) 0%;
    }
  }
</style>
