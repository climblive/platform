import { expect, test, vi } from "vitest";
import { categorySchema } from "./category";
import { competitionSchema } from "./competition";
import { competitorSchema } from "./competitor";
import { rulesUpdatedEventSchema } from "./event";
import { problemSchema } from "./problem";
import { raffleSchema } from "./raffle";
import { scoreEngineDescriptorSchema } from "./scoreEngine";

vi.mock("@climblive/lib/utils", () => import("../utils/zod"));

const id = "019a85f0-0000-7000-8000-000000000001";
const rules = {
  qualifyingProblems: 10,
  finalists: 5,
  usePoints: true,
  pooledPoints: false,
  maxAttempts: 0,
  pointDeduction: 0,
};

const cases = [
  {
    name: "competition",
    schema: competitionSchema,
    key: "id",
    value: {
      ...rules,
      id,
      ownership: { organizerId: 1 },
      country: "SE",
      name: "Competition",
      gracePeriod: 0,
      nameRetentionTime: 0,
      created: new Date(),
      registeredCompetitors: 0,
    },
  },
  {
    name: "category",
    schema: categorySchema,
    value: {
      id: 1,
      competitionId: id,
      name: "Open",
      timeBegin: new Date(),
      timeEnd: new Date(),
    },
  },
  {
    name: "competitor",
    schema: competitorSchema,
    value: {
      id: 1,
      competitionId: id,
      registrationCode: "ABCD1234",
      withdrawnFromFinals: false,
      disqualified: false,
    },
  },
  {
    name: "problem",
    schema: problemSchema,
    value: {
      id: 1,
      competitionId: id,
      number: 1,
      holdColorPrimary: "#ffffff",
      pointsTop: 100,
      zone1Enabled: false,
      zone2Enabled: false,
    },
  },
  { name: "raffle", schema: raffleSchema, value: { id: 1, competitionId: id } },
  {
    name: "rules event",
    schema: rulesUpdatedEventSchema,
    value: { ...rules, competitionId: id },
  },
  {
    name: "score engine",
    schema: scoreEngineDescriptorSchema,
    value: { competitionId: id, instanceId: id },
  },
];

for (const { name, schema, value, key = "competitionId" } of cases) {
  test(`${name} accepts UUID competition IDs`, () => {
    expect(schema.safeParse(value).success).toBe(true);
  });

  test(`${name} rejects numeric or malformed competition IDs`, () => {
    for (const invalid of [1, "1", "not-a-uuid", ""]) {
      expect(schema.safeParse({ ...value, [key]: invalid }).success).toBe(
        false,
      );
    }
  });
}

test("other entity identifiers remain numeric", () => {
  expect(raffleSchema.safeParse({ id, competitionId: id }).success).toBe(false);
});
