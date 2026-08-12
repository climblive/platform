import type { Problem, Tick } from "@climblive/lib/models";
import { SvelteMap, SvelteSet } from "svelte/reactivity";

type Feature = "zone1" | "zone2" | "top";

export class TickBuilder {
  #problemId: number;
  #features: Feature[];
  #reachedFeatures: SvelteSet<Feature>;
  #attempts: SvelteMap<Feature, number>;
  #tick: Omit<Tick, "id" | "timestamp">;

  constructor(problem: Problem, tick?: Tick) {
    this.#problemId = problem.id;

    this.#features = [];
    this.#features.push("zone1");
    this.#features.push("zone2");
    this.#features.push("top");

    this.#reachedFeatures = new SvelteSet();

    for (const feature of this.#features) {
      switch (feature) {
        case "zone1":
          if (tick?.zone1) {
            this.#reachedFeatures.add("zone1");
          }
          break;
        case "zone2":
          if (tick?.zone2) {
            this.#reachedFeatures.add("zone2");
          }
          break;
        case "top":
          if (tick?.top) {
            this.#reachedFeatures.add("top");
          }
          break;
      }
    }

    this.#attempts = new SvelteMap();
    this.#attempts.set("zone1", tick?.attemptsZone1 ?? 0);
    this.#attempts.set("zone2", tick?.attemptsZone2 ?? 0);
    this.#attempts.set("top", tick?.attemptsTop ?? 0);

    this.#tick = $derived<Omit<Tick, "id" | "timestamp">>({
      problemId: this.#problemId,
      zone1: this.#reachedFeatures.has("zone1"),
      attemptsZone1: this.#attempts.get("zone1") ?? 0,
      zone2: this.#reachedFeatures.has("zone2"),
      attemptsZone2: this.#attempts.get("zone2") ?? 0,
      top: this.#reachedFeatures.has("top"),
      attemptsTop: this.#attempts.get("top") ?? 0,
    });
  }

  public get tick(): Omit<Tick, "id" | "timestamp"> {
    return this.#tick;
  }

  public addAttempt(): void {
    if (!this.canAddAttempt()) {
      return;
    }

    if (!this.#reachedFeatures.has("top")) {
      this.#attempts.set("top", (this.#attempts.get("top") ?? 0) + 1);
    }

    if (!this.#reachedFeatures.has("zone2")) {
      this.#attempts.set("zone2", (this.#attempts.get("zone2") ?? 0) + 1);
    }

    if (!this.#reachedFeatures.has("zone1")) {
      this.#attempts.set("zone1", (this.#attempts.get("zone1") ?? 0) + 1);
    }
  }

  public subtractAttempt(): void {
    if (!this.canSubtractAttempt()) {
      return;
    }

    if (!this.#reachedFeatures.has("top")) {
      this.#attempts.set("top", (this.#attempts.get("top") ?? 0) - 1);
    }

    if (!this.#reachedFeatures.has("zone2")) {
      this.#attempts.set("zone2", (this.#attempts.get("zone2") ?? 0) - 1);
    }

    if (!this.#reachedFeatures.has("zone1")) {
      this.#attempts.set("zone1", (this.#attempts.get("zone1") ?? 0) - 1);
    }
  }

  public canAddAttempt(): boolean {
    if (this.#reachedFeatures.has("top")) {
      return false;
    }

    for (const attempts of this.#attempts.values()) {
      if (attempts === 999) {
        return false;
      }
    }

    return true;
  }

  public canSubtractAttempt(): boolean {
    if (this.#reachedFeatures.has("top")) {
      return false;
    }

    for (const attempts of this.#attempts.values()) {
      if (attempts === 0) {
        return false;
      }
    }

    for (let k = this.#features.length - 1; k >= 0; k--) {
      const f0 = this.#features[k - 1];
      const f1 = this.#features[k];

      const a0 = this.#attempts.get(f0);
      const a1 = this.#attempts.get(f1);

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

  public normalizeAttempts() {
    if (this.#reachedFeatures.size === 0) {
      const firstFeature = this.#features[0];

      if (!firstFeature) {
        return;
      }

      const attemptsFirstFeature = this.#attempts.get(firstFeature) ?? 0;

      for (const feature of this.#features) {
        this.#attempts.set(feature, attemptsFirstFeature);
      }
    }

    for (let k = 0; k < this.#features.length; k++) {
      const f1 = this.#features[k];
      const f2 = this.#features[k + 1];
      const f3 = this.#features[k + 2];

      if (f2 === undefined || f2 === undefined) {
        continue;
      }

      if (
        this.#reachedFeatures.has(f1) &&
        !this.#reachedFeatures.has(f2) &&
        !this.#reachedFeatures.has(f3)
      ) {
        for (let j = k + 2; j < this.#features.length; j++) {
          const feature = this.#features[j];
          this.#attempts.set(feature, this.#attempts.get(f2) ?? 0);
        }

        return;
      }
    }
  }

  public reachFeature(feature: Feature): void {
    this.addAttempt();

    switch (feature) {
      case "top":
        this.#reachedFeatures.add("top");
      // eslint-disable-next-line no-fallthrough
      case "zone2":
        this.#reachedFeatures.add("zone2");
      // eslint-disable-next-line no-fallthrough
      case "zone1":
        this.#reachedFeatures.add("zone1");
    }
  }

  public unreachFeature(feature: Feature): void {
    switch (feature) {
      case "zone1":
        this.#reachedFeatures.delete("zone1");
      // eslint-disable-next-line no-fallthrough
      case "zone2":
        this.#reachedFeatures.delete("zone2");
      // eslint-disable-next-line no-fallthrough
      case "top":
        this.#reachedFeatures.delete("top");
    }

    this.normalizeAttempts();
  }
}
