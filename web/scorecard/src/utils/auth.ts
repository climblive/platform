import { scorecardSessionSchema, type ScorecardSession } from "@/types";
import { ApiClient, ContenderCredentialsProvider } from "@climblive/lib";
import type { Contender } from "@climblive/lib/models";
import { z } from "@climblive/lib/utils";
import type { QueryClient } from "@tanstack/svelte-query";
import { add } from "date-fns";
import type { Writable } from "svelte/store";

export const authenticateContender = async (
  code: string,
  queryClient: QueryClient,
  session: Writable<ScorecardSession>,
): Promise<Contender> => {
  const contender = await queryClient.fetchQuery({
    queryKey: ["contender", { code }],
    queryFn: async () => ApiClient.getInstance().findContender(code),
  });

  const contest = await queryClient.fetchQuery({
    queryKey: ["contest", { id: contender.contestId }],
    queryFn: async () =>
      ApiClient.getInstance().getContest(contender.contestId),
  });

  const provider = new ContenderCredentialsProvider(code);
  ApiClient.getInstance().setCredentialsProvider(provider);

  session.update((current) => {
    const now = new Date();
    const contestEndTime = contest.timeEnd || now;
    const baseTime = new Date(
      Math.max(contestEndTime.getTime(), now.getTime()),
    );
    const expiryTime = add(baseTime, { hours: 12 });

    const updatedSession: ScorecardSession = {
      ...current,
      contenderId: contender.id,
      contestId: contender.contestId,
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
    ["contender", { id: contender.id }],
    () => contender,
  );

  return contender;
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
