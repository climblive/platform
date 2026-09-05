import type { Problem, Tick } from "@climblive/lib/models";
import { SvelteMap } from "svelte/reactivity";

export type Feature = "zone1" | "zone2" | "top";

const allFeatures: Feature[] = ["zone1", "zone2", "top"];

export class TickBuilder {
  #problemId: number;
  #features: Feature[];
  #reachedFeatures: SvelteMap<Feature, number>;
  #attempts: number;

  constructor(
    problemId: number,
    features: Feature[],
    attempts: number,
    reachedFeatures: Map<Feature, number>,
  ) {
    this.#problemId = problemId;
    this.#features = [...features];
    this.#attempts = $state(attempts);
    this.#reachedFeatures = new SvelteMap(reachedFeatures);

    const attemptsTop = this.#reachedFeatures.get("top");
    if (attemptsTop !== undefined && attemptsTop !== this.#attempts) {
      throw new Error("Attempts mismatch");
    }
  }

  static from(
    problem: Pick<Problem, "id" | "zone1Enabled" | "zone2Enabled">,
    tick?: Pick<
      Tick,
      | "top"
      | "attemptsTop"
      | "zone2"
      | "attemptsZone2"
      | "zone1"
      | "attemptsZone1"
    >,
  ): TickBuilder {
    const features: Feature[] = [];
    if (problem.zone1Enabled) {
      features.push("zone1");
    }
    if (problem.zone2Enabled) {
      features.push("zone2");
    }
    features.push("top");

    const reachedFeatures = new SvelteMap<Feature, number>();

    if (tick) {
      for (const feature of features) {
        switch (feature) {
          case "zone1":
            if (tick.zone1) {
              reachedFeatures.set("zone1", tick.attemptsZone1);
            }
            break;
          case "zone2":
            if (tick.zone2) {
              reachedFeatures.set("zone2", tick.attemptsZone2);
            }
            break;
          case "top":
            if (tick.top) {
              reachedFeatures.set("top", tick.attemptsTop);
            }
            break;
        }
      }
    }

    return new TickBuilder(
      problem.id,
      features,
      tick?.attemptsTop ?? 0,
      reachedFeatures,
    );
  }

  public get problemId(): number {
    return this.#problemId;
  }

  public get features(): readonly Feature[] {
    return this.#features;
  }

  public get reachedFeatures(): ReadonlyMap<Feature, number> {
    return this.#reachedFeatures;
  }

  public get attempts(): number {
    return this.#attempts;
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

    if (this.#attempts >= 999) {
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

export function buildTick(
  builder: Pick<TickBuilder, "problemId" | "attempts" | "reachedFeatures">,
): Omit<Tick, "id" | "timestamp" | "revision"> {
  const hasReached = (feature: Feature): boolean => {
    let featureReached = false;

    for (let k = allFeatures.length - 1; k >= 0; k--) {
      const f = allFeatures[k];
      if (builder.reachedFeatures.has(f)) {
        featureReached = true;
      }

      if (f === feature) {
        return featureReached;
      }
    }

    return false;
  };

  const calculateImplicitAttempts = (feature: Feature): number => {
    let attempts = builder.attempts;

    for (let k = allFeatures.length - 1; k >= 0; k--) {
      const f = allFeatures[k];
      const a = builder.reachedFeatures.get(f);
      if (a !== undefined) {
        attempts = a;
      }

      if (f === feature) {
        return attempts;
      }
    }

    return attempts;
  };

  return {
    problemId: builder.problemId,
    zone1: hasReached("zone1"),
    attemptsZone1: calculateImplicitAttempts("zone1"),
    zone2: hasReached("zone2"),
    attemptsZone2: calculateImplicitAttempts("zone2"),
    top: hasReached("top"),
    attemptsTop: calculateImplicitAttempts("top"),
  };
}
