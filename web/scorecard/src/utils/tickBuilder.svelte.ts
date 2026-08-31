import type { Problem, Tick } from "@climblive/lib/models";
import { SvelteMap } from "svelte/reactivity";

type Feature = "zone1" | "zone2" | "top";

const allFeatures: Feature[] = ["zone1", "zone2", "top"];

export class TickBuilder {
  #problemId: number;
  #features: Feature[];
  #reachedFeatures: SvelteMap<Feature, number>;
  #attempts: number;
  #tick: Omit<Tick, "id" | "timestamp">;

  constructor(problem: Problem, tick?: Tick) {
    this.#problemId = problem.id;

    this.#features = [];
    if (problem.zone1Enabled) {
      this.#features.push("zone1");
    }
    if (problem.zone2Enabled) {
      this.#features.push("zone2");
    }
    this.#features.push("top");

    this.#reachedFeatures = new SvelteMap();
    this.#attempts = $state(0);

    if (tick) {
      for (const feature of this.#features) {
        switch (feature) {
          case "zone1":
            if (tick.zone1) {
              this.#reachedFeatures.set("zone1", tick.attemptsZone1);
            }
            break;
          case "zone2":
            if (tick.zone2) {
              this.#reachedFeatures.set("zone2", tick.attemptsZone2);
            }
            break;
          case "top":
            if (tick.top) {
              this.#reachedFeatures.set("top", tick.attemptsTop);
            }
            break;
        }
      }

      this.#attempts = tick.attemptsTop;
    }

    this.#tick = $derived<Omit<Tick, "id" | "timestamp">>({
      problemId: this.#problemId,
      zone1: this.#hasReached("zone1"),
      attemptsZone1: this.#calculateImplicitAttempts("zone1"),
      zone2: this.#hasReached("zone2"),
      attemptsZone2: this.#calculateImplicitAttempts("zone2"),
      top: this.#hasReached("top"),
      attemptsTop: this.#calculateImplicitAttempts("top"),
    });
  }

  #hasReached(feature: Feature): boolean {
    let featureReached = false;

    for (let k = allFeatures.length - 1; k >= 0; k--) {
      const f = allFeatures[k];
      if (this.#reachedFeatures.has(f)) {
        featureReached = true;
      }

      if (f === feature) {
        return featureReached;
      }
    }

    return false;
  }

  #calculateImplicitAttempts(feature: Feature): number {
    let attempts = this.#attempts;

    for (let k = allFeatures.length - 1; k >= 0; k--) {
      const f = allFeatures[k];
      const a = this.#reachedFeatures.get(f);
      if (a !== undefined) {
        attempts = a;
      }

      if (f === feature) {
        return attempts;
      }
    }

    return attempts;
  }

  public get tick(): Omit<Tick, "id" | "timestamp"> {
    return this.#tick;
  }

  public addAttempt(): void {
    if (!this.canAddAttempt()) {
      return;
    }

    this.#attempts += 1;
  }

  public subtractAttempt(): void {
    if (!this.canSubtractAttempt()) {
      return;
    }

    this.#attempts -= 1;
  }

  public canAddAttempt(): boolean {
    if (this.#reachedFeatures.has("top")) {
      return false;
    }

    if (this.#attempts === 999) {
      return false;
    }

    return true;
  }

  public canSubtractAttempt(): boolean {
    if (this.#reachedFeatures.has("top")) {
      return false;
    }

    if (this.#attempts === 0) {
      return false;
    }

    return this.#attempts > Math.max(...this.#reachedFeatures.values());
  }

  public reachFeature(feature: Feature): void {
    this.addAttempt();

    const featureIndex = this.#features.findIndex((f) => f === feature);
    if (featureIndex === -1) {
      return;
    }

    for (let k = featureIndex; k >= 0; k--) {
      const f = this.#features[k];

      if (!this.#reachedFeatures.has(f)) {
        this.#reachedFeatures.set(f, this.#attempts);
      }
    }
  }

  public unreachFeature(feature: Feature): void {
    const featureIndex = this.#features.findIndex((f) => f === feature);
    if (featureIndex === -1) {
      return;
    }

    for (let k = featureIndex; k < this.#features.length; k++) {
      const f = this.#features[k];

      this.#reachedFeatures.delete(f);
    }

    this.subtractAttempt();
  }
}
