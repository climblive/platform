import type { Problem, Tick } from "@climblive/lib/models";
import { SvelteMap } from "svelte/reactivity";

type Feature = "zone1" | "zone2" | "top";

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
      attemptsZone1: this.#calculateAttempts("zone1"),
      zone2: this.#hasReached("zone2"),
      attemptsZone2: this.#calculateAttempts("zone2"),
      top: this.#hasReached("top"),
      attemptsTop: this.#calculateAttempts("top"),
    });
  }

  #hasReached(feature: Feature): boolean {
    let featureReached = false;

    for (let k = this.#features.length - 1; k >= 0; k--) {
      const f = this.#features[k];
      if (this.#reachedFeatures.has(f)) {
        featureReached = true;
      }

      if (f === feature) {
        return featureReached;
      }
    }

    return false;
  }

  #calculateAttempts(feature: Feature): number {
    const attempts = this.#reachedFeatures.get(feature);
    if (attempts !== undefined) {
      return attempts;
    }

    return this.#attempts;
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
    const topFeature = this.#features[this.#features.length - 1];

    if (this.#reachedFeatures.has(topFeature)) {
      return false;
    }

    if (this.#attempts === 999) {
      return false;
    }

    return true;
  }

  public canSubtractAttempt(): boolean {
    const topFeature = this.#features[this.#features.length - 1];

    if (this.#reachedFeatures.has(topFeature)) {
      return false;
    }

    if (this.#attempts === 0) {
      return false;
    }

    for (let k = this.#features.length - 1; k >= 0; k--) {
      const f0 = this.#features[k - 1];
      const f1 = this.#features[k];

      const a0 = this.#reachedFeatures.get(f0);
      const a1 = this.#reachedFeatures.get(f1);

      if (
        a0 !== undefined &&
        a1 !== undefined &&
        this.#reachedFeatures.has(f0) &&
        a0 === a1
      ) {
        return false;
      }
    }

    return true;
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
