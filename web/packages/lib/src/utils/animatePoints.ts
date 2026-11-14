const property = "--climblive-animated-points";
const duration = 1_000;
let supported: boolean | undefined;

const supportsAnimation = () => {
  if (supported !== undefined) {
    return supported;
  }

  supported = false;
  let probe: HTMLSpanElement | undefined;
  let animation: Animation | undefined;

  try {
    if (
      typeof CSS?.registerProperty !== "function" ||
      typeof Element.prototype.animate !== "function" ||
      !CSS.supports("transition", `${property} 1s`)
    ) {
      return false;
    }

    CSS.registerProperty({
      name: property,
      syntax: "<number>",
      inherits: false,
      initialValue: "0",
    });

    probe = document.createElement("span");
    probe.style.cssText =
      "position:fixed;visibility:hidden;pointer-events:none";
    document.body.append(probe);
    animation = probe.animate(
      [{ [property]: "10.5" }, { [property]: "30.5" }, { [property]: "10.5" }],
      { duration: 200, fill: "both" },
    );
    animation.pause();

    const target = probe;
    const interpolation = animation;
    supported = [0, 25, 50, 75, 100, 125, 150, 175, 200].every((time) => {
      interpolation.currentTime = time;
      const value = Number(getComputedStyle(target).getPropertyValue(property));
      const expected = 10.5 + Math.min(time, 200 - time) / 5;
      return Math.abs(value - expected) < 0.001;
    });
  } catch {
    supported = false;
  } finally {
    animation?.cancel();
    probe?.remove();
  }

  return supported;
};

export const animatePoints = (
  element: HTMLElement,
  from: number,
  to: number,
  update: (value: number | undefined) => void,
) => {
  let animation: Animation | undefined;
  let frame: number | undefined;
  let timeout: ReturnType<typeof setTimeout> | undefined;
  let motion: MediaQueryList | undefined;
  let stopped = false;

  const stop = () => {
    if (stopped) {
      return;
    }
    stopped = true;
    if (frame !== undefined) {
      cancelAnimationFrame(frame);
    }
    clearTimeout(timeout);
    animation?.cancel();
    motion?.removeEventListener?.("change", stop);
    document.removeEventListener("visibilitychange", stop);
    update(undefined);
  };

  try {
    motion = matchMedia("(prefers-reduced-motion: reduce)");
    if (
      motion.matches ||
      typeof motion.addEventListener !== "function" ||
      document.hidden ||
      !Number.isFinite(from) ||
      !Number.isFinite(to) ||
      Math.max(Math.abs(from), Math.abs(to)) > 2_147_483_647 ||
      from === to ||
      !supportsAnimation()
    ) {
      return stop;
    }

    animation = element.animate(
      [{ [property]: String(from) }, { [property]: String(to) }],
      { duration, easing: "ease-out" },
    );
    const started = performance.now();

    const render = () => {
      try {
        if (
          performance.now() - started >= duration ||
          animation?.playState !== "running"
        ) {
          stop();
          return;
        }

        const raw = getComputedStyle(element).getPropertyValue(property).trim();
        const value = raw === "" ? NaN : Number(raw);
        if (
          !Number.isFinite(value) ||
          value < Math.min(from, to) ||
          value > Math.max(from, to)
        ) {
          supported = false;
          stop();
          return;
        }

        update(value);
        frame = requestAnimationFrame(render);
      } catch {
        supported = false;
        stop();
      }
    };

    motion.addEventListener("change", stop);
    document.addEventListener("visibilitychange", stop);
    animation.onfinish = stop;
    timeout = setTimeout(stop, duration);
    frame = requestAnimationFrame(render);
  } catch {
    stop();
  }

  return stop;
};
