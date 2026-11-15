import { expect, test } from "vitest";
import { contenderSchema } from "./contender";
import { contenderResultUpdatedEventSchema } from "./event";
import { resultSchema } from "./result";
import { scoreboardEntrySchema } from "./scoreboard";

const result = {
  timestamp: "2026-11-15T12:00:00Z",
  contenderId: 42,
  score: "123.4p",
  placement: 2,
  finalist: true,
  rankOrder: 1,
};

const parsedResult = { ...result, timestamp: new Date(result.timestamp) };

test("parses ranked results and result update events", () => {
  expect(resultSchema.parse(result)).toEqual(parsedResult);
  expect(contenderResultUpdatedEventSchema.parse(result)).toEqual(parsedResult);
});

test.each([
  {
    name: "contender",
    schema: contenderSchema,
    input: {
      id: 42,
      competitionId: 1,
      registrationCode: "ABCD0042",
      withdrawnFromFinals: false,
      disqualified: false,
    },
  },
  {
    name: "scoreboard entry",
    schema: scoreboardEntrySchema,
    input: {
      contenderId: 42,
      categoryId: 1,
      name: "Climber",
      withdrawnFromFinals: false,
      disqualified: false,
    },
  },
])("parses optional results on a $name", ({ schema, input }) => {
  expect(schema.parse({ ...input, result }).result).toEqual(parsedResult);
  expect(schema.parse(input).result).toBeUndefined();
  expect(
    schema.safeParse({ ...input, result: { ...result, placement: "2" } })
      .success,
  ).toBe(false);
});
