import { SvelteDate } from "svelte/reactivity";

export class SyncedTime {
  #interval: number;
  #time: Date;
  #tickCb: ((time: Date) => void) | undefined;
  #intervalTimerId: number | undefined;
  #enabled: boolean;

  constructor(interval: number = 1_000, tickCb?: (time: Date) => void) {
    this.#interval = interval;
    // eslint-disable-next-line svelte/prefer-svelte-reactivity
    this.#time = $state(new Date());
    this.#tickCb = tickCb;
    this.#enabled = false;
  }

  start() {
    this.#enabled = true;

    this.tick();
  }

  stop() {
    this.#enabled = false;

    if (this.#intervalTimerId !== undefined) {
      clearTimeout(this.#intervalTimerId);
    }
  }

  public get current() {
    return this.#time;
  }

  private tick = () => {
    if (!this.#enabled) {
      return;
    }

    this.#time = new SvelteDate();
    this.#tickCb?.(this.#time);

    const firefoxEarlyWakeUpCompensation = 1;

    const drift = Date.now() % this.#interval;
    const next = this.#interval - drift + firefoxEarlyWakeUpCompensation;

    this.#intervalTimerId = setTimeout(this.tick, next);
  };
}
