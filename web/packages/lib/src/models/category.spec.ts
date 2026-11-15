import { describe, expect, it } from "vitest";
import { categorySchema } from "./category";
import { contenderSchema } from "./contender";
import { contenderPublicInfoUpdatedEventSchema } from "./event";
import { scoreboardEntrySchema } from "./scoreboard";

describe("category schemas", () => {
  it("parses categories and their competition times", () => {
    expect(
      categorySchema.parse({
        id: 42,
        competitionId: 123,
        name: "Open",
        timeBegin: "2026-11-15T10:00:00Z",
        timeEnd: "2026-11-15T16:00:00Z",
      }),
    ).toEqual({
      id: 42,
      competitionId: 123,
      name: "Open",
      timeBegin: new Date("2026-11-15T10:00:00Z"),
      timeEnd: new Date("2026-11-15T16:00:00Z"),
    });
  });

  it.each([
    { name: "contender", schema: contenderSchema },
    { name: "scoreboard entry", schema: scoreboardEntrySchema },
    {
      name: "public info event",
      schema: contenderPublicInfoUpdatedEventSchema,
    },
  ])("preserves categoryId in a $name", ({ schema }) => {
    const parsed = schema.parse({
      id: 1,
      contenderId: 1,
      competitionId: 123,
      categoryId: 42,
      registrationCode: "ABCD1234",
      name: "Climber",
      withdrawnFromFinals: false,
      disqualified: false,
    });

    expect(parsed.categoryId).toBe(42);
  });
});
