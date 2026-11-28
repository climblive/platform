import { expect, test } from "vitest";
import { competitionSchema } from "./competition";

const competition = {
  id: 123,
  ownership: { organizerId: 1 },
  country: "SE",
  name: "Competition",
  qualifyingProblems: 10,
  finalists: 6,
  usePoints: true,
  pooledPoints: false,
  maxAttempts: 0,
  pointDeduction: 0,
  gracePeriod: 0,
  nameRetentionTime: 0,
  created: "2026-11-27T18:00:00Z",
  registeredCompetitors: 0,
};

test.each([
  { duplicatedBy: null, duplicatedFrom: null },
  { duplicatedBy: "2026-11-27T18:00:00Z", duplicatedFrom: 42 },
  { duplicatedBy: "2026-11-27T18:00:00Z", duplicatedFrom: null },
  { duplicatedBy: null, duplicatedFrom: 42 },
])("parses nullable duplication metadata %j", (metadata) => {
  const parsed = competitionSchema.parse({ ...competition, ...metadata });

  expect(parsed.duplicatedFrom).toBe(metadata.duplicatedFrom);
  expect(parsed.duplicatedBy).toEqual(
    metadata.duplicatedBy === null ? null : new Date(metadata.duplicatedBy),
  );
});
