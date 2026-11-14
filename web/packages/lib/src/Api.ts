import type { AxiosInstance, RawAxiosRequestHeaders } from "axios";
import axios from "axios";
import { z } from "zod/v4";
import {
  competitionSchema,
  scoreboardEntrySchema,
  scoreEngineDescriptorSchema,
  type CompClassPatch,
  type CompClassTemplate,
  type ContenderPatch,
  type CompetitionID,
  type CompetitionPatch,
  type CompetitionTemplate,
  type OrganizerInviteID,
  type OrganizerPatch,
  type OrganizerTemplate,
  type ProblemPatch,
  type ProblemTemplate,
  type ScoreEngineInstanceID,
  type Tick,
} from "./models";
import { compClassSchema } from "./models/compClass";
import { contenderSchema } from "./models/contender";
import { serviceStatusSchema } from "./models/health";
import { organizerSchema } from "./models/organizer";
import { organizerInviteSchema } from "./models/organizerInvite";
import { pointValueSchema } from "./models/pointValue";
import { problemSchema } from "./models/problem";
import { raffleSchema, raffleWinnerSchema } from "./models/raffle";
import type {
  CreateContendersArguments,
  StartScoreEngineArguments,
} from "./models/rest";
import { tickSchema } from "./models/tick";
import { userSchema } from "./models/user";

interface ApiCredentialsProvider {
  getAuthHeaders(): RawAxiosRequestHeaders;
}

export class ContenderCredentialsProvider implements ApiCredentialsProvider {
  private registrationCode: string;

  constructor(registrationCode: string) {
    this.registrationCode = registrationCode;
  }

  getAuthHeaders = (): RawAxiosRequestHeaders => {
    const headers: RawAxiosRequestHeaders = {};

    headers["Authorization"] = `Regcode ${this.registrationCode}`;

    return headers;
  };
}

export class OrganizerCredentialsProvider implements ApiCredentialsProvider {
  private jwt: string;

  constructor(jwt: string) {
    this.jwt = jwt;
  }

  getAuthHeaders = (): RawAxiosRequestHeaders => {
    const headers: RawAxiosRequestHeaders = {};

    headers["Authorization"] = `Bearer ${this.jwt}`;

    return headers;
  };
}

export class ApiClient {
  private static instance: ApiClient;
  private axiosInstance: AxiosInstance;
  private credentialsProvider: ApiCredentialsProvider | undefined;

  private constructor() {
    this.axiosInstance = axios.create({
      baseURL: "/api",
      timeout: 10_000,
    });
  }

  public static getInstance(): ApiClient {
    if (!ApiClient.instance) {
      ApiClient.instance = new ApiClient();
    }

    return ApiClient.instance;
  }

  setCredentialsProvider = (credentialsProvider: ApiCredentialsProvider) => {
    this.credentialsProvider = credentialsProvider;
  };

  getSelf = async () => {
    const endpoint = "/users/self";

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return userSchema.parse(result.data);
  };

  findContender = async (registrationCode: string) => {
    const endpoint = `/codes/${registrationCode}/contender`;

    const result = await this.axiosInstance.get(endpoint);

    return contenderSchema.parse(result.data);
  };

  getContender = async (id: number) => {
    const endpoint = `/contenders/${id}`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return contenderSchema.parse(result.data);
  };

  getPointValuesByContender = async (id: number) => {
    const endpoint = `/contenders/${id}/points`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(pointValueSchema).parse(result.data);
  };

  getContendersByCompetition = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/contenders`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(contenderSchema).parse(result.data);
  };

  patchContender = async (id: number, patch: ContenderPatch) => {
    const endpoint = `/contenders/${id}`;

    const result = await this.axiosInstance.patch(endpoint, patch, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return contenderSchema.parse(result.data);
  };

  scrubContender = async (id: number) => {
    const endpoint = `/contenders/${id}/scrub`;

    const result = await this.axiosInstance.post(endpoint, null, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return contenderSchema.parse(result.data);
  };

  createContenders = async (
    competitionId: number,
    args: CreateContendersArguments,
  ) => {
    const endpoint = `/competitions/${competitionId}/contenders`;

    const result = await this.axiosInstance.post(endpoint, args, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(contenderSchema).parse(result.data);
  };

  getCompetition = async (id: number) => {
    const endpoint = `/competitions/${id}`;

    const result = await this.axiosInstance.get(endpoint);

    return competitionSchema.parse(result.data);
  };

  getAllCompetitions = async () => {
    const endpoint = `/competitions`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(competitionSchema).parse(result.data);
  };

  createCompetition = async (
    organizerId: number,
    template: CompetitionTemplate,
  ) => {
    const endpoint = `/organizers/${organizerId}/competitions`;

    const result = await this.axiosInstance.post(endpoint, template, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return competitionSchema.parse(result.data);
  };

  patchCompetition = async (id: number, patch: CompetitionPatch) => {
    const endpoint = `/competitions/${id}`;

    const result = await this.axiosInstance.patch(endpoint, patch, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return competitionSchema.parse(result.data);
  };

  duplicateCompetition = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/duplicate`;

    const result = await this.axiosInstance.post(endpoint, undefined, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return competitionSchema.parse(result.data);
  };

  archiveCompetition = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/archive`;

    const result = await this.axiosInstance.post(endpoint, undefined, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return competitionSchema.parse(result.data);
  };

  restoreCompetition = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/restore`;

    const result = await this.axiosInstance.post(endpoint, undefined, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return competitionSchema.parse(result.data);
  };

  transferCompetition = async (
    competitionId: number,
    newOrganizerId: number,
  ) => {
    const endpoint = `/competitions/${competitionId}/transfer`;

    const result = await this.axiosInstance.post(
      endpoint,
      { newOrganizerId },
      {
        headers: this.credentialsProvider?.getAuthHeaders(),
      },
    );

    return competitionSchema.parse(result.data);
  };

  getCompetitionsByOrganizer = async (organizerId: number) => {
    const endpoint = `/organizers/${organizerId}/competitions`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(competitionSchema).parse(result.data);
  };

  getProblem = async (problemId: number) => {
    const endpoint = `/problems/${problemId}`;

    const result = await this.axiosInstance.get(endpoint);

    return problemSchema.parse(result.data);
  };

  getProblemsByCompetition = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/problems`;

    const result = await this.axiosInstance.get(endpoint);

    return z.array(problemSchema).parse(result.data);
  };

  createProblem = async (competitionId: number, template: ProblemTemplate) => {
    const endpoint = `/competitions/${competitionId}/problems`;

    const result = await this.axiosInstance.post(endpoint, template, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return problemSchema.parse(result.data);
  };

  patchProblem = async (id: number, patch: ProblemPatch) => {
    const endpoint = `/problems/${id}`;

    const result = await this.axiosInstance.patch(endpoint, patch, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return problemSchema.parse(result.data);
  };

  deleteProblem = async (id: number) => {
    const endpoint = `/problems/${id}`;

    await this.axiosInstance.delete(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });
  };

  getCompClass = async (compClassId: number) => {
    const endpoint = `/comp-classes/${compClassId}`;

    const result = await this.axiosInstance.get(endpoint);

    return compClassSchema.parse(result.data);
  };

  getCompClasses = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/comp-classes`;

    const result = await this.axiosInstance.get(endpoint);

    return z.array(compClassSchema).parse(result.data);
  };

  createCompClass = async (
    competitionId: number,
    template: CompClassTemplate,
  ) => {
    const endpoint = `/competitions/${competitionId}/comp-classes`;

    const result = await this.axiosInstance.post(endpoint, template, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return compClassSchema.parse(result.data);
  };

  deleteCompClass = async (id: number) => {
    const endpoint = `/comp-classes/${id}`;

    await this.axiosInstance.delete(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });
  };

  patchCompClass = async (id: number, patch: CompClassPatch) => {
    const endpoint = `/comp-classes/${id}`;

    const result = await this.axiosInstance.patch(endpoint, patch, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return compClassSchema.parse(result.data);
  };

  getTicksByContender = async (contenderId: number) => {
    const endpoint = `/contenders/${contenderId}/ticks`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(tickSchema).parse(result.data);
  };

  getTicksByCompetition = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/ticks`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(tickSchema).parse(result.data);
  };

  putTick = async (
    contenderId: number,
    tick: Omit<Tick, "id" | "timestamp">,
  ) => {
    const endpoint = `/contenders/${contenderId}/ticks`;

    const result = await this.axiosInstance.put(endpoint, tick, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return tickSchema.parse(result.data);
  };

  deleteTick = async (tickId: number) => {
    const endpoint = `/ticks/${tickId}`;

    await this.axiosInstance.delete(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });
  };

  getScoreboard = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/scoreboard`;

    const result = await this.axiosInstance.get(endpoint);

    return z.array(scoreboardEntrySchema).parse(result.data);
  };

  getScoreEnginesByCompetition = async (competitionId: CompetitionID) => {
    const endpoint = `/competitions/${competitionId}/score-engines`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(scoreEngineDescriptorSchema).parse(result.data);
  };

  getScoreEngines = async () => {
    const endpoint = "/score-engines";

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(scoreEngineDescriptorSchema).parse(result.data);
  };

  startScoreEngine = async (
    competitionId: CompetitionID,
    args: StartScoreEngineArguments,
  ) => {
    const endpoint = `/competitions/${competitionId}/score-engines`;

    const result = await this.axiosInstance.post(endpoint, args, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.string().uuid().parse(result.data);
  };

  stopScoreEngine = async (instanceId: ScoreEngineInstanceID) => {
    const endpoint = `/score-engines/${instanceId}`;

    await this.axiosInstance.delete(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });
  };

  getRaffle = async (raffleId: number) => {
    const endpoint = `/raffles/${raffleId}`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return raffleSchema.parse(result.data);
  };

  getRaffles = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/raffles`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(raffleSchema).parse(result.data);
  };

  createRaffle = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/raffles`;

    const result = await this.axiosInstance.post(endpoint, undefined, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return raffleSchema.parse(result.data);
  };

  deleteRaffle = async (raffleId: number) => {
    const endpoint = `/raffles/${raffleId}`;

    await this.axiosInstance.delete(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });
  };

  drawRaffleWinner = async (raffleId: number) => {
    const endpoint = `/raffles/${raffleId}/winners`;

    const result = await this.axiosInstance.post(endpoint, undefined, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return raffleWinnerSchema.parse(result.data);
  };

  getRaffleWinners = async (raffleId: number) => {
    const endpoint = `/raffles/${raffleId}/winners`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(raffleWinnerSchema).parse(result.data);
  };

  createOrganizer = async (template: OrganizerTemplate) => {
    const endpoint = "/organizers";

    const result = await this.axiosInstance.post(endpoint, template, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return organizerSchema.parse(result.data);
  };

  getOrganizer = async (organizerId: number) => {
    const endpoint = `/organizers/${organizerId}`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return organizerSchema.parse(result.data);
  };

  patchOrganizer = async (organizerId: number, patch: OrganizerPatch) => {
    const endpoint = `/organizers/${organizerId}`;

    const result = await this.axiosInstance.patch(endpoint, patch, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return organizerSchema.parse(result.data);
  };

  getOrganizerInvites = async (organizerId: number) => {
    const endpoint = `/organizers/${organizerId}/invites`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(organizerInviteSchema).parse(result.data);
  };

  getOrganizerInvite = async (inviteId: OrganizerInviteID) => {
    const endpoint = `/invites/${inviteId}`;

    const result = await this.axiosInstance.get(endpoint, {});

    return organizerInviteSchema.parse(result.data);
  };

  deleteOrganizerInvite = async (inviteId: OrganizerInviteID) => {
    const endpoint = `/invites/${inviteId}`;

    await this.axiosInstance.delete(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });
  };

  acceptOrganizerInvite = async (inviteId: OrganizerInviteID) => {
    const endpoint = `/invites/${inviteId}/accept`;

    await this.axiosInstance.post(endpoint, undefined, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });
  };

  createOrganizerInvite = async (organizerId: number) => {
    const endpoint = `/organizers/${organizerId}/invites`;

    const result = await this.axiosInstance.post(endpoint, undefined, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return organizerInviteSchema.parse(result.data);
  };

  getUsersByOrganizer = async (organizerId: number) => {
    const endpoint = `/organizers/${organizerId}/users`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(userSchema).parse(result.data);
  };

  downloadResults = async (competitionId: number) => {
    const endpoint = `/competitions/${competitionId}/results`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
      responseType: "blob",
    });

    return result.data;
  };

  getHealth = async () => {
    const endpoint = "/health";

    const result = await this.axiosInstance.get(endpoint, {});

    return z.array(serviceStatusSchema).parse(result.data);
  };

  getVersion = async () => {
    const endpoint = "/version";

    const result = await this.axiosInstance.get(endpoint, {});

    return z.string().parse(result.data);
  };
}
