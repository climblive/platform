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

  it("restores a saved session", () => {
    const getItem = vi.fn().mockReturnValue(
      JSON.stringify([
        {
          contenderId: 7,
          competitionId: 42,
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

  it("does not restore expired sessions", () => {
    vi.stubGlobal("localStorage", {
      getItem: () =>
        JSON.stringify([
          {
            contenderId: 7,
            competitionId: 42,
            registrationCode: "ABCD0001",
            expiryTime: "2026-11-14T11:00:00Z",
          },
        ]),
    });

    expect(readStoredSessions()).toEqual([]);
  });

  it("does not restore legacy sessions", () => {
    vi.stubGlobal("localStorage", {
      getItem: () =>
        JSON.stringify([
          {
            contenderId: 7,
            contestId: 42,
            registrationCode: "ABCD0001",
            expiryTime: "2026-11-15T00:00:00Z",
          },
        ]),
    });

    expect(readStoredSessions()).toEqual([]);
  });
});
