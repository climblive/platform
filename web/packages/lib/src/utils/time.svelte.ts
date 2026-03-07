export class SyncedTime {
  #interval: number;
  #time: Date;
  #intervalTimerId: number | undefined;
  #ticking: boolean;

  constructor(interval: number = 1_000) {
    this.#interval = interval;
    // eslint-disable-next-line svelte/prefer-svelte-reactivity
    this.#time = $state(new Date());
    this.#ticking = false;
  }

  start() {
    this.#ticking = true;

    this.tick();
  }

  stop() {
    this.#ticking = false;

    if (this.#intervalTimerId !== undefined) {
      clearTimeout(this.#intervalTimerId);
      this.#intervalTimerId = undefined;
    }
  }

  public get current() {
    return this.#time;
  }

  private tick = () => {
    if (!this.#ticking) {
      return;
    }

    // eslint-disable-next-line svelte/prefer-svelte-reactivity
    const now = new Date();

    this.#time = now;

    const firefoxEarlyWakeUpCompensation = 1;

    const drift = now.getTime() % this.#interval;
    const next = this.#interval - drift + firefoxEarlyWakeUpCompensation;

    this.#intervalTimerId = setTimeout(this.tick, next);
  };
}
