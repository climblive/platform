import { describe, expect, it } from "vitest";
import { buildTick, TickMutator, type Feature } from "./tickMutator.svelte";

const PROBLEM_ID = 123;
const ALL_FEATURES: Feature[] = ["zone1", "zone2", "top"];
const FLASH = new Map<Feature, number>([
  ["zone1", 1],
  ["zone2", 1],
  ["top", 1],
]);
const NO_RESULT = new Map<Feature, number>();

describe(buildTick.name, () => {
  it("should build a tick with implicit features and attempts", () => {
    expect(
      buildTick(
        PROBLEM_ID,
        new TickMutator(ALL_FEATURES, 4, new Map([["zone2", 3]])),
      ),
    ).toEqual({
      problemId: PROBLEM_ID,
      zone1: true,
      attemptsZone1: 3,
      zone2: true,
      attemptsZone2: 3,
      top: false,
      attemptsTop: 4,
    });
  });
});

describe(TickMutator.name, () => {
  it("should expose its state", () => {
    const mutator = new TickMutator(ALL_FEATURES, 1, FLASH);

    expect(mutator.features).toEqual(ALL_FEATURES);
    expect(mutator.attempts).toEqual(1);
    expect(mutator.reachedFeatures).toEqual(FLASH);
  });

  it("should not allow attempts to differ from attempts required for top", () => {
    expect(() => new TickMutator(ALL_FEATURES, 2, FLASH)).toThrow();
  });

  describe(TickMutator.from.name, () => {
    it("should create an empty mutator", () => {
      const mutator = TickMutator.from({
        zone1Enabled: true,
        zone2Enabled: true,
      });

      expect(mutator.features).toEqual(ALL_FEATURES);
      expect(mutator.attempts).toEqual(0);
      expect(mutator.reachedFeatures).toEqual(NO_RESULT);
    });

    it("should restore an existing tick", () => {
      const mutator = TickMutator.from(
        { zone1Enabled: true, zone2Enabled: true },
        {
          zone1: true,
          attemptsZone1: 1,
          zone2: true,
          attemptsZone2: 2,
          top: true,
          attemptsTop: 3,
        },
      );

      expect(mutator.features).toEqual(ALL_FEATURES);
      expect(mutator.attempts).toEqual(3);
      expect(mutator.reachedFeatures).toEqual(
        new Map<Feature, number>([
          ["zone1", 1],
          ["zone2", 2],
          ["top", 3],
        ]),
      );
    });

    it("should ignore reached features that are not enabled", () => {
      const mutator = TickMutator.from(
        { zone1Enabled: false, zone2Enabled: false },
        {
          zone1: true,
          attemptsZone1: 1,
          zone2: true,
          attemptsZone2: 2,
          top: true,
          attemptsTop: 3,
        },
      );

      expect(mutator.features).toEqual(["top"]);
      expect(mutator.attempts).toEqual(3);
      expect(mutator.reachedFeatures).toEqual(
        new Map<Feature, number>([["top", 3]]),
      );
    });
  });

  describe(TickMutator.prototype.canAddAttempt.name, () => {
    it("should return true if another attempt can be added", () => {
      const mutator = new TickMutator(ALL_FEATURES, 998, NO_RESULT);

      expect(mutator.canAddAttempt()).toEqual(true);
    });

    it("should return false if top is reached", () => {
      const mutator = new TickMutator(ALL_FEATURES, 1, FLASH);

      expect(mutator.canAddAttempt()).toEqual(false);
    });

    it("should return false if attempts is 999 or more", () => {
      const mutator = new TickMutator(ALL_FEATURES, 999, NO_RESULT);

      expect(mutator.canAddAttempt()).toEqual(false);
    });
  });

  describe(TickMutator.prototype.canSubtractAttempt.name, () => {
    it("should return true if the latest attempt reached no features", () => {
      const mutator = new TickMutator(ALL_FEATURES, 2, new Map([["zone1", 1]]));

      expect(mutator.canSubtractAttempt()).toEqual(true);
    });

    it("should return false if top is reached", () => {
      const mutator = new TickMutator(ALL_FEATURES, 1, FLASH);

      expect(mutator.canSubtractAttempt()).toEqual(false);
    });

    it("should return false if there are no attempts", () => {
      const mutator = new TickMutator(ALL_FEATURES, 0, NO_RESULT);

      expect(mutator.canSubtractAttempt()).toEqual(false);
    });

    it("should return false if the latest attempt reached a feature", () => {
      const mutator = new TickMutator(ALL_FEATURES, 1, new Map([["zone1", 1]]));

      expect(mutator.canSubtractAttempt()).toEqual(false);
    });
  });

  describe(TickMutator.prototype.addAttempt.name, () => {
    it("should add an attempt", () => {
      const mutator = new TickMutator(ALL_FEATURES, 0, NO_RESULT);

      const attemptAdded = mutator.addAttempt();

      expect(attemptAdded).toEqual(true);
      expect(mutator.attempts).toEqual(1);
    });

    it("should not add an attempt if another attempt cannot be added", () => {
      const mutator = new TickMutator(ALL_FEATURES, 1, FLASH);

      const attemptAdded = mutator.addAttempt();

      expect(mutator.attempts).toEqual(1);
      expect(attemptAdded).toEqual(false);
    });
  });

  describe(TickMutator.prototype.subtractAttempt.name, () => {
    it("should subtract an attempt", () => {
      const mutator = new TickMutator(ALL_FEATURES, 1, NO_RESULT);

      const attemptSubtracted = mutator.subtractAttempt();

      expect(mutator.attempts).toEqual(0);
      expect(attemptSubtracted).toEqual(true);
    });

    it("should not subtract an attempt if the latest attempt reached a feature", () => {
      const mutator = new TickMutator(ALL_FEATURES, 1, new Map([["zone1", 1]]));

      const attemptSubtracted = mutator.subtractAttempt();

      expect(mutator.attempts).toEqual(1);
      expect(attemptSubtracted).toEqual(false);
    });
  });

  describe(TickMutator.prototype.reachFeature.name, () => {
    it("should reach all preceding features when reaching top", () => {
      const mutator = new TickMutator(ALL_FEATURES, 1, NO_RESULT);

      mutator.reachFeature("top");

      expect(mutator.attempts).toEqual(2);
      expect(mutator.reachedFeatures).toEqual(
        new Map<Feature, number>([
          ["zone1", 2],
          ["zone2", 2],
          ["top", 2],
        ]),
      );
    });

    it("should reach the feature and all preceding features", () => {
      const mutator = new TickMutator(ALL_FEATURES, 1, NO_RESULT);

      mutator.reachFeature("zone2");

      expect(mutator.attempts).toEqual(2);
      expect(mutator.reachedFeatures).toEqual(
        new Map<Feature, number>([
          ["zone1", 2],
          ["zone2", 2],
        ]),
      );
    });

    it("should not reach a feature that is not enabled", () => {
      const mutator = new TickMutator(["top"], 0, NO_RESULT);

      mutator.reachFeature("zone1");

      expect(mutator.attempts).toEqual(1);
      expect(mutator.reachedFeatures).toEqual(NO_RESULT);
    });

    it("should not reach a feature if no more attempts can be added", () => {
      const mutator = new TickMutator(ALL_FEATURES, 999, NO_RESULT);

      mutator.reachFeature("top");

      expect(mutator.attempts).toEqual(999);
      expect(mutator.reachedFeatures).toEqual(NO_RESULT);
    });
  });

  describe(TickMutator.prototype.unreachFeature.name, () => {
    it("should unreach the feature and all subsequent features", () => {
      const mutator = new TickMutator(
        ALL_FEATURES,
        3,
        new Map<Feature, number>([
          ["zone1", 1],
          ["zone2", 2],
          ["top", 3],
        ]),
      );

      mutator.unreachFeature("zone2");

      expect(mutator.attempts).toEqual(2);
      expect(mutator.reachedFeatures).toEqual(
        new Map<Feature, number>([["zone1", 1]]),
      );
    });

    it("should not change a feature that is not enabled", () => {
      const mutator = new TickMutator(["top"], 1, NO_RESULT);

      mutator.unreachFeature("zone1");

      expect(mutator.attempts).toEqual(1);
      expect(mutator.reachedFeatures).toEqual(NO_RESULT);
    });
  });
});
