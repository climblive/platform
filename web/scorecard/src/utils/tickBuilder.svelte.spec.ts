import { describe, expect, it } from "vitest";
import { TickBuilder, type Feature } from "./tickBuilder.svelte";

const PROBLEM_ID = 123;
const ALL_FEATURES: Feature[] = ["zone1", "zone2", "top"];
const FLASH = new Map<Feature, number>([
  ["zone1", 1],
  ["zone2", 1],
  ["top", 1],
]);
const NO_LUCK = new Map<Feature, number>();

describe(TickBuilder.name, () => {
  it("should derive the tick from reached features", () => {
    const builder = new TickBuilder(
      PROBLEM_ID,
      ["zone2", "top"],
      4,
      new Map<Feature, number>([["zone2", 3]]),
    );

    expect(builder.tick).toEqual({
      problemId: PROBLEM_ID,
      zone1: true,
      attemptsZone1: 3,
      zone2: true,
      attemptsZone2: 3,
      top: false,
      attemptsTop: 4,
    });
  });

  it("should not allow attempts to differ from attempts required for top", () => {
    expect(
      () => new TickBuilder(PROBLEM_ID, ALL_FEATURES, 2, FLASH),
    ).toThrowError("Attempts mismatch");
  });

  describe(TickBuilder.from.name, () => {
    it("should create an empty tick", () => {
      const builder = TickBuilder.from({
        id: PROBLEM_ID,
        zone1Enabled: true,
        zone2Enabled: true,
      });

      expect(builder.tick).toEqual({
        problemId: PROBLEM_ID,
        zone1: false,
        attemptsZone1: 0,
        zone2: false,
        attemptsZone2: 0,
        top: false,
        attemptsTop: 0,
      });
    });

    it("should restore an existing tick", () => {
      const builder = TickBuilder.from(
        { id: PROBLEM_ID, zone1Enabled: true, zone2Enabled: true },
        {
          zone1: true,
          attemptsZone1: 1,
          zone2: true,
          attemptsZone2: 2,
          top: true,
          attemptsTop: 3,
        },
      );

      expect(builder.tick).toEqual({
        problemId: PROBLEM_ID,
        zone1: true,
        attemptsZone1: 1,
        zone2: true,
        attemptsZone2: 2,
        top: true,
        attemptsTop: 3,
      });
    });

    it("should ignore reached features that are not enabled", () => {
      const builder = TickBuilder.from(
        { id: PROBLEM_ID, zone1Enabled: false, zone2Enabled: false },
        {
          zone1: true,
          attemptsZone1: 1,
          zone2: true,
          attemptsZone2: 2,
          top: false,
          attemptsTop: 3,
        },
      );

      expect(builder.tick).toEqual({
        problemId: PROBLEM_ID,
        zone1: false,
        attemptsZone1: 3,
        zone2: false,
        attemptsZone2: 3,
        top: false,
        attemptsTop: 3,
      });
    });
  });

  describe(TickBuilder.prototype.canAddAttempt.name, () => {
    it("should return true if another attempt can be added", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 998, NO_LUCK);

      expect(builder.canAddAttempt()).toEqual(true);
    });

    it("should return false if top is reached", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 1, FLASH);

      expect(builder.canAddAttempt()).toEqual(false);
    });

    it("should return false if attempts is 999 or more", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 999, NO_LUCK);

      expect(builder.canAddAttempt()).toEqual(false);
    });
  });

  describe(TickBuilder.prototype.canSubtractAttempt.name, () => {
    it("should return true if the latest attempt reached no features", () => {
      const builder = new TickBuilder(
        PROBLEM_ID,
        ALL_FEATURES,
        2,
        new Map([["zone1", 1]]),
      );

      expect(builder.canSubtractAttempt()).toEqual(true);
    });

    it("should return false if top is reached", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 1, FLASH);

      expect(builder.canSubtractAttempt()).toEqual(false);
    });

    it("should return false if there are no attempts", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 0, NO_LUCK);

      expect(builder.canSubtractAttempt()).toEqual(false);
    });

    it("should return false if the latest attempt reached a feature", () => {
      const builder = new TickBuilder(
        PROBLEM_ID,
        ALL_FEATURES,
        1,
        new Map([["zone1", 1]]),
      );

      expect(builder.canSubtractAttempt()).toEqual(false);
    });
  });

  describe(TickBuilder.prototype.addAttempt.name, () => {
    it("should add an attempt", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 0, NO_LUCK);

      builder.addAttempt();

      expect(builder.tick.attemptsTop).toEqual(1);
      expect(builder.tick.attemptsZone2).toEqual(1);
      expect(builder.tick.attemptsZone1).toEqual(1);
    });

    it("should not add an attempt if another attempt cannot be added", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 1, FLASH);

      builder.addAttempt();

      expect(builder.tick.attemptsTop).toEqual(1);
    });
  });

  describe(TickBuilder.prototype.subtractAttempt.name, () => {
    it("should subtract an attempt", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 1, NO_LUCK);

      builder.subtractAttempt();

      expect(builder.tick.attemptsTop).toEqual(0);
      expect(builder.tick.attemptsZone2).toEqual(0);
      expect(builder.tick.attemptsZone1).toEqual(0);
    });

    it("should not subtract an attempt if the latest attempt reached a feature", () => {
      const builder = new TickBuilder(
        PROBLEM_ID,
        ALL_FEATURES,
        1,
        new Map([["zone1", 1]]),
      );

      builder.subtractAttempt();

      expect(builder.tick.attemptsTop).toEqual(1);
    });
  });

  describe(TickBuilder.prototype.reachFeature.name, () => {
    it("should reach the feature and all preceding features", () => {
      const builder = new TickBuilder(PROBLEM_ID, ALL_FEATURES, 1, NO_LUCK);

      builder.reachFeature("zone2");

      expect(builder.tick).toEqual({
        problemId: PROBLEM_ID,
        zone1: true,
        attemptsZone1: 2,
        zone2: true,
        attemptsZone2: 2,
        top: false,
        attemptsTop: 2,
      });
    });

    it("should not reach a feature that is not enabled", () => {
      const builder = new TickBuilder(PROBLEM_ID, ["top"], 0, NO_LUCK);

      builder.reachFeature("zone1");

      expect(builder.tick).toEqual({
        problemId: PROBLEM_ID,
        zone1: false,
        attemptsZone1: 1,
        zone2: false,
        attemptsZone2: 1,
        top: false,
        attemptsTop: 1,
      });
    });
  });

  describe(TickBuilder.prototype.unreachFeature.name, () => {
    it("should unreach the feature and all subsequent features", () => {
      const builder = new TickBuilder(
        PROBLEM_ID,
        ALL_FEATURES,
        3,
        new Map<Feature, number>([
          ["zone1", 1],
          ["zone2", 2],
          ["top", 3],
        ]),
      );

      builder.unreachFeature("zone2");

      expect(builder.tick).toEqual({
        problemId: PROBLEM_ID,
        zone1: true,
        attemptsZone1: 1,
        zone2: false,
        attemptsZone2: 2,
        top: false,
        attemptsTop: 2,
      });
    });

    it("should not change a feature that is not enabled", () => {
      const builder = new TickBuilder(PROBLEM_ID, ["top"], 1, NO_LUCK);

      builder.unreachFeature("zone1");

      expect(builder.tick).toEqual({
        problemId: PROBLEM_ID,
        zone1: false,
        attemptsZone1: 1,
        zone2: false,
        attemptsZone2: 1,
        top: false,
        attemptsTop: 1,
      });
    });
  });
});
