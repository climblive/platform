import {
  createMutation,
  createQuery,
  QueryClient,
  useQueryClient,
  type QueryKey,
} from "@tanstack/svelte-query";
import { ApiClient } from "../Api";
import {
  type Competition,
  type CompetitionPatch,
  type CompetitionTemplate,
} from "../models";

export const getCompetitionQuery = (competitionId: string) =>
  createQuery(() => ({
    queryKey: ["competition", { id: competitionId }],
    queryFn: async () => ApiClient.getInstance().getCompetition(competitionId),
  }));

export const getAllCompetitionsQuery = (
  options?: Partial<Parameters<typeof createQuery<Competition[]>>[0]>,
) =>
  createQuery(() => ({
    ...options,
    queryKey: ["competitions"],
    queryFn: async () => ApiClient.getInstance().getAllCompetitions(),
  }));

export const getCompetitionsByOrganizerQuery = (
  organizerId: number,
  options?: Partial<Parameters<typeof createQuery<Competition[]>>[0]>,
) =>
  createQuery(() => ({
    ...options,
    queryKey: ["competitions", { organizerId }],
    queryFn: async () =>
      ApiClient.getInstance().getCompetitionsByOrganizer(organizerId),
  }));

export const createCompetitionMutation = (organizerId: number) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (template: CompetitionTemplate) =>
      ApiClient.getInstance().createCompetition(organizerId, template),
    onSuccess: (newCompetition) => {
      let queryKey: QueryKey = ["competitions", { organizerId }];

      client.setQueryData<Competition[]>(queryKey, (oldCompetitions) => {
        return [...(oldCompetitions ?? []), newCompetition];
      });

      queryKey = ["competition", { id: newCompetition.id }];

      client.setQueryData<Competition>(queryKey, newCompetition);
    },
  }));
};

export const patchCompetitionMutation = (competitionId: string) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (patch: CompetitionPatch) =>
      ApiClient.getInstance().patchCompetition(competitionId, patch),
    onSuccess: (patchedCompetition) => {
      let queryKey: QueryKey = [
        "competitions",
        { organizerId: patchedCompetition.ownership.organizerId },
      ];

      client.setQueryData<Competition[]>(queryKey, (oldCompetitions) => {
        if (oldCompetitions === undefined) {
          return undefined;
        }

        return oldCompetitions.map((competition) => {
          if (competition.id === patchedCompetition.id) {
            return patchedCompetition;
          }

          return competition;
        });
      });

      queryKey = ["competition", { id: competitionId }];

      client.setQueryData<Competition>(queryKey, patchedCompetition);
    },
  }));
};

export const duplicateCompetitionMutation = (competitionId: string) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: () =>
      ApiClient.getInstance().duplicateCompetition(competitionId),
    onSuccess: (duplicatedCompetition) => {
      let queryKey: QueryKey = [
        "competitions",
        { organizerId: duplicatedCompetition.ownership.organizerId },
      ];

      client.setQueryData<Competition[]>(queryKey, (oldCompetitions) => {
        return [...(oldCompetitions ?? []), duplicatedCompetition];
      });

      queryKey = ["competition", { id: duplicatedCompetition.id }];

      client.setQueryData<Competition>(queryKey, duplicatedCompetition);
    },
  }));
};

export const archiveCompetitionMutation = (competitionId: string) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => ApiClient.getInstance().archiveCompetition(competitionId),
    onSuccess: (archivedCompetition) => {
      let queryKey: QueryKey = ["competitions"];

      client.setQueriesData<Competition[]>(
        { queryKey, exact: false },
        (oldCompetitions) => {
          if (oldCompetitions === undefined) {
            return undefined;
          }

          return oldCompetitions.map((competition) => {
            if (competition.id === archivedCompetition.id) {
              return archivedCompetition;
            }

            return competition;
          });
        },
      );

      queryKey = ["competition", { id: competitionId }];

      client.setQueryData<Competition>(queryKey, archivedCompetition);
    },
  }));
};

export const restoreCompetitionMutation = (competitionId: string) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => ApiClient.getInstance().restoreCompetition(competitionId),
    onSuccess: (restoredCompetition) => {
      let queryKey: QueryKey = ["competitions"];

      client.setQueriesData<Competition[]>(
        { queryKey, exact: false },
        (oldCompetitions) => {
          if (oldCompetitions === undefined) {
            return undefined;
          }

          return oldCompetitions.map((competition) => {
            if (competition.id === restoredCompetition.id) {
              return restoredCompetition;
            }

            return competition;
          });
        },
      );

      queryKey = ["competition", { id: competitionId }];

      client.setQueryData<Competition>(queryKey, restoredCompetition);
    },
  }));
};

export const transferCompetitionMutation = (competitionId: string) => {
  const client = useQueryClient();

  return createMutation(() => ({
    mutationFn: (newOrganizerId: number) =>
      ApiClient.getInstance().transferCompetition(
        competitionId,
        newOrganizerId,
      ),
    onSuccess: (transferredCompetition, newOrganizerId) => {
      const oldOrganizerId = client.getQueryData<Competition>([
        "competition",
        { id: competitionId },
      ])?.ownership.organizerId;

      if (oldOrganizerId) {
        client.setQueryData<Competition[]>(
          ["competitions", { organizerId: oldOrganizerId }],
          (competitions) => {
            if (competitions === undefined) {
              return undefined;
            }

            return competitions.filter(({ id }) => id !== competitionId);
          },
        );
      }

      client.setQueryData<Competition[]>(
        ["competitions", { organizerId: newOrganizerId }],
        (competitions) => {
          return [...(competitions ?? []), transferredCompetition];
        },
      );

      client.setQueryData<Competition>(
        ["competition", { id: competitionId }],
        transferredCompetition,
      );

      client.invalidateQueries({ queryKey: ["competitions"] });
    },
  }));
};

export const updateRulesInQueryCache = (
  queryClient: QueryClient,
  competitionId: string,
  updatedRules: Pick<
    Competition,
    | "qualifyingProblems"
    | "finalists"
    | "usePoints"
    | "pooledPoints"
    | "maxAttempts"
    | "pointDeduction"
  >,
) => {
  let queryKey: QueryKey = ["competitions"];

  queryClient.setQueriesData<Competition[]>(
    { queryKey, exact: false },
    (oldCompetitions) => {
      if (oldCompetitions === undefined) {
        return undefined;
      }

      return oldCompetitions.map((competition) => {
        if (competition.id === competitionId) {
          return { ...competition, ...updatedRules };
        }

        return competition;
      });
    },
  );

  queryKey = ["competition", { id: competitionId }];

  queryClient.setQueryData<Competition>(queryKey, (oldCompetition) => {
    if (!oldCompetition) {
      return undefined;
    }

    return { ...oldCompetition, ...updatedRules };
  });
};
