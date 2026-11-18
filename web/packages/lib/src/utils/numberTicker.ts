let supported: boolean | undefined;

export const supportsNumberTicker = () => {
  if (supported !== undefined) {
    return supported;
  }

  supported = false;
  let animation: Animation | undefined;
  let sampler: Animation | undefined;

  try {
    if (
      typeof IntersectionObserver !== "function" ||
      typeof Animation !== "function" ||
      typeof KeyframeEffect !== "function" ||
      typeof Element.prototype.animate !== "function" ||
      typeof requestAnimationFrame !== "function" ||
      typeof matchMedia !== "function"
    ) {
      return false;
    }

    const options = { duration: 100, fill: "both" as const };
    const effect = new KeyframeEffect(null, null, options);
    sampler = new Animation(effect, null);
    animation = document.createElement("span").animate(null, options);
    animation.pause();

    const player = animation;
    const probe = sampler;
    supported = [0, 25, 50, 75, 100].every((time) => {
      probe.currentTime = time;
      player.currentTime = time;

      return [effect, player.effect].every((effect) => {
        const progress = effect?.getComputedTiming().progress;
        return progress != null && Math.abs(progress - time / 100) < 0.001;
      });
    });
  } catch {
    supported = false;
  } finally {
    animation?.cancel();
    sampler?.cancel();
  }

  return supported;
};
