import { describe, expect, it } from "vitest";
import { TickBuilder } from "./tickBuilder.svelte";

describe(TickBuilder.name, () => {
  describe(TickBuilder.prototype.canAddAttempt.name, () => {
    it("should return false if top is reached", () => {
      const builder = new TickBuilder(
        { id: 1, zone1Enabled: false, zone2Enabled: false },
        {
          top: true,
          attemptsTop: 1,
          zone2: true,
          attemptsZone2: 1,
          zone1: true,
          attemptsZone1: 1,
        },
      );

      expect(builder.canAddAttempt()).toEqual(false);
    });

    it("should return false if attempts is 999 or more", () => {
      const builder = new TickBuilder(
        { id: 1, zone1Enabled: false, zone2Enabled: false },
        {
          top: false,
          attemptsTop: 999,
          zone2: false,
          attemptsZone2: 999,
          zone1: false,
          attemptsZone1: 999,
        },
      );

      expect(builder.canAddAttempt()).toEqual(false);
    });
  });

  describe(TickBuilder.prototype.canSubtractAttempt.name, () => {
    it("should return false if top is reached", () => {
      const builder = new TickBuilder(
        { id: 1, zone1Enabled: false, zone2Enabled: false },
        {
          top: true,
          attemptsTop: 1,
          zone2: true,
          attemptsZone2: 1,
          zone1: true,
          attemptsZone1: 1,
        },
      );

      expect(builder.canSubtractAttempt()).toEqual(false);
    });

    it("should return false if attempts is 0", () => {
      const builder = new TickBuilder(
        { id: 1, zone1Enabled: false, zone2Enabled: false },
        {
          top: false,
          attemptsTop: 0,
          zone2: false,
          attemptsZone2: 0,
          zone1: false,
          attemptsZone1: 0,
        },
      );

      expect(builder.canSubtractAttempt()).toEqual(false);
    });
  });

  describe(TickBuilder.prototype.addAttempt.name, () => {
    it("should add attempt to empty tick", () => {
      const builder = new TickBuilder(
        { id: 1, zone1Enabled: false, zone2Enabled: false },
        undefined,
      );

      builder.addAttempt();

      expect(builder.tick.attemptsTop).toEqual(1);
      expect(builder.tick.attemptsZone2).toEqual(1);
      expect(builder.tick.attemptsZone1).toEqual(1);
    });
  });
});
