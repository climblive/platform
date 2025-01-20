import { scorecardSessionSchema, type ScorecardSession } from "@/types";
import { ApiClient, ContenderCredentialsProvider } from "@climblive/lib";
import type { Contender } from "@climblive/lib/models";
import type { QueryClient } from "@tanstack/svelte-query";
import { differenceInHours } from "date-fns";
import type { Writable } from "svelte/store";
import * as z from "zod";

export const authenticateContender = async (
  code: string,
  queryClient: QueryClient,
  session: Writable<ScorecardSession>,
): Promise<Contender> => {
  const contender = await ApiClient.getInstance().findContender(code);

  const provider = new ContenderCredentialsProvider(code);
  ApiClient.getInstance().setCredentialsProvider(provider);

  session.update((current) => {
    const updatedSession: ScorecardSession = {
      ...current,
      contenderId: contender.id,
      contestId: contender.contestId,
      registrationCode: code,
      timestamp: new Date(),
    };

    let sessions = readStoredSessions();
    sessions = sessions.filter(session => session.registrationCode !== updatedSession.registrationCode)
    sessions.push(updatedSession);

    localStorage.setItem("session", JSON.stringify(sessions));

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

  const data = localStorage.getItem("session");
  if (data) {
    try {
      const obj = JSON.parse(data);
      const storedSessions = z.array(scorecardSessionSchema).parse(obj);

      for (const storedSession of storedSessions) {
        if (differenceInHours(new Date(), storedSession.timestamp) < 12) {
          sessions.push(storedSession);
        }
      }
    } catch {
      /* discard corrupt session data */
    }
  }

  sessions.sort((s1: ScorecardSession, s2: ScorecardSession) => { return s2.timestamp.getTime() - s1.timestamp.getTime() })

  return sessions.slice(0, 3);
}