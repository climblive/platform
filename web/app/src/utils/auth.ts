import type { ScorecardSession } from "@/types";
import { ApiClient, ContenderCredentialsProvider } from "@climblive/lib";
import type { Contender } from "@climblive/lib/models";
import type { QueryClient } from "@tanstack/svelte-query";
import type { Writable } from "svelte/store";

export const authenticateContender = async (
  code: string,
  queryClient: QueryClient,
  session: Writable<ScorecardSession>,
): Promise<Contender> => {
  const provider = new ContenderCredentialsProvider(code);
  ApiClient.getInstance().setCredentialsProvider(provider);

  const contender = await ApiClient.getInstance().findContender(code);

  session.update((current) => {
    const updatedSession = {
      ...current,
      contenderId: contender.id,
      contestId: contender.contestId,
      registrationCode: code,
      timestamp: new Date(),
    };

    localStorage.setItem("session", JSON.stringify(updatedSession));

    return updatedSession;
  });

  queryClient.setQueryData(
    ["contender", { id: contender.id }],
    () => contender,
  );

  return contender;
};
