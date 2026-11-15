import { scorecardSessionSchema, type ScorecardSession } from "@/types";
import { ApiClient, CompetitorCredentialsProvider } from "@climblive/lib";
import type { Competitor } from "@climblive/lib/models";
import { z } from "@climblive/lib/utils";
import type { QueryClient } from "@tanstack/svelte-query";
import { add } from "date-fns";
import type { Writable } from "svelte/store";

export const authenticateCompetitor = async (
  code: string,
  queryClient: QueryClient,
  session: Writable<ScorecardSession>,
): Promise<Competitor> => {
  const competitor = await queryClient.fetchQuery({
    queryKey: ["competitor", { code }],
    queryFn: async () => ApiClient.getInstance().findCompetitor(code),
  });

  const competition = await queryClient.fetchQuery({
    queryKey: ["competition", { id: competitor.competitionId }],
    queryFn: async () =>
      ApiClient.getInstance().getCompetition(competitor.competitionId),
  });

  const provider = new CompetitorCredentialsProvider(code);
  ApiClient.getInstance().setCredentialsProvider(provider);

  session.update((current) => {
    const now = new Date();
    const competitionEndTime = competition.timeEnd || now;
    const baseTime = new Date(
      Math.max(competitionEndTime.getTime(), now.getTime()),
    );
    const expiryTime = add(baseTime, { hours: 12 });

    const updatedSession: ScorecardSession = {
      ...current,
      competitorId: competitor.id,
      competitionId: competitor.competitionId,
      registrationCode: code,
      expiryTime,
    };

    let sessions = readStoredSessions();
    const predicate = ({ registrationCode }: ScorecardSession) =>
      registrationCode !== updatedSession.registrationCode;
    sessions = sessions.filter(predicate);
    sessions.splice(0, 0, updatedSession);

    localStorage.setItem("sessions", JSON.stringify(sessions));

    return updatedSession;
  });

  queryClient.setQueryData(
    ["competitor", { id: competitor.id }],
    () => competitor,
  );

  return competitor;
};

export const readStoredSessions = (): ScorecardSession[] => {
  const sessions: ScorecardSession[] = [];

  const data = localStorage.getItem("sessions");
  if (data) {
    try {
      const obj = JSON.parse(data);
      const storedSessions = z.array(scorecardSessionSchema).parse(obj);

      const now = new Date();
      for (const storedSession of storedSessions) {
        if (new Date(storedSession.expiryTime) > now) {
          sessions.push(storedSession);
        }
      }
    } catch {
      /* discard corrupt session data */
    }
  }

  sessions.sort((s1: ScorecardSession, s2: ScorecardSession) => {
    return s2.expiryTime.getTime() - s1.expiryTime.getTime();
  });

  return sessions.slice(0, 3);
};
