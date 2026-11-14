import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { readStoredSessions } from "./auth";

describe(readStoredSessions.name, () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-11-14T12:00:00Z"));
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
  });

  it.each([
    { competitionId: 42 },
    { contestId: 42 },
    { competitionId: 42, contestId: 99 },
  ])("restores a saved session with %j", (competition) => {
    const getItem = vi.fn().mockReturnValue(
      JSON.stringify([
        {
          contenderId: 7,
          ...competition,
          registrationCode: "ABCD0001",
          expiryTime: "2026-11-15T00:00:00Z",
        },
      ]),
    );
    vi.stubGlobal("localStorage", { getItem });

    expect(readStoredSessions()).toEqual([
      {
        contenderId: 7,
        competitionId: 42,
        registrationCode: "ABCD0001",
        expiryTime: new Date("2026-11-15T00:00:00Z"),
      },
    ]);
    expect(getItem).toHaveBeenCalledWith("sessions");
  });

  it("does not restore expired legacy sessions", () => {
    vi.stubGlobal("localStorage", {
      getItem: () =>
        JSON.stringify([
          {
            contenderId: 7,
            contestId: 42,
            registrationCode: "ABCD0001",
            expiryTime: "2026-11-14T11:00:00Z",
          },
        ]),
    });

    expect(readStoredSessions()).toEqual([]);
  });
});
