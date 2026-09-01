import { describe, expect, it } from "vitest";
import { TickBuilder, type Feature } from "./tickBuilder.svelte";

const ALL_FEATURES: Feature[] = ["top", "zone2", "zone1"] as const;
const FLASH = new Map<Feature, number>([
  ["top", 1],
  ["zone2", 1],
  ["zone1", 1],
]);

const NO_LUCK = new Map<Feature, number>();

describe(TickBuilder.name, () => {
  describe(TickBuilder.prototype.canAddAttempt.name, () => {
    it("should return false if top is reached", () => {
      const builder = new TickBuilder(123, ALL_FEATURES, 1, FLASH);

      expect(builder.canAddAttempt()).toEqual(false);
    });

    it("should return false if attempts is 999 or more", () => {
      const builder = new TickBuilder(123, ALL_FEATURES, 999, NO_LUCK);

      expect(builder.canAddAttempt()).toEqual(false);
    });
  });

  describe(TickBuilder.prototype.canSubtractAttempt.name, () => {
    it("should return false if top is reached", () => {
      const builder = new TickBuilder(123, ALL_FEATURES, 1, FLASH);

      expect(builder.canSubtractAttempt()).toEqual(false);
    });

    it("should return false if attempts is 0", () => {
      const builder = new TickBuilder(123, ALL_FEATURES, 0, NO_LUCK);

      expect(builder.canSubtractAttempt()).toEqual(false);
    });
  });

  describe(TickBuilder.prototype.addAttempt.name, () => {
    it("should add attempt to empty tick", () => {
      const builder = new TickBuilder(123, ALL_FEATURES, 0, NO_LUCK);

      builder.addAttempt();

      expect(builder.tick.attemptsTop).toEqual(1);
      expect(builder.tick.attemptsZone2).toEqual(1);
      expect(builder.tick.attemptsZone1).toEqual(1);
    });
  });
});
