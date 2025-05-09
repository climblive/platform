import type { AxiosInstance, RawAxiosRequestHeaders } from "axios";
import axios from "axios";
import { z } from "zod";
import {
  contestSchema,
  scoreboardEntrySchema,
  type CompClassTemplate,
  type ContenderPatch,
  type ContestID,
  type ContestTemplate,
  type ProblemPatch,
  type ProblemTemplate,
  type ScoreEngineInstanceID,
  type Tick,
} from "./models";
import { compClassSchema } from "./models/compClass";
import { contenderSchema } from "./models/contender";
import { problemSchema } from "./models/problem";
import type {
  CreateContendersArguments,
  StartScoreEngineArguments,
} from "./models/rest";
import { tickSchema } from "./models/tick";
import { getApiUrl } from "./utils/config";

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
      baseURL: getApiUrl(),
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

  getContendersByContest = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/contenders`;

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

  createContenders = async (
    contestId: number,
    args: CreateContendersArguments,
  ) => {
    const endpoint = `/contests/${contestId}/contenders`;

    const result = await this.axiosInstance.post(endpoint, args, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(contenderSchema).parse(result.data);
  };

  getContest = async (id: number) => {
    const endpoint = `/contests/${id}`;

    const result = await this.axiosInstance.get(endpoint);

    return contestSchema.parse(result.data);
  };

  createContest = async (organizerId: number, template: ContestTemplate) => {
    const endpoint = `/organizers/${organizerId}/contests`;

    const result = await this.axiosInstance.post(endpoint, template, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return contestSchema.parse(result.data);
  };

  getContestsByOrganizer = async (organizerId: number) => {
    const endpoint = `/organizers/${organizerId}/contests`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(contestSchema).parse(result.data);
  };

  getProblem = async (problemId: number) => {
    const endpoint = `/problems/${problemId}`;

    const result = await this.axiosInstance.get(endpoint);

    return problemSchema.parse(result.data);
  };

  getProblems = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/problems`;

    const result = await this.axiosInstance.get(endpoint);

    return z.array(problemSchema).parse(result.data);
  };

  createProblem = async (contestId: number, template: ProblemTemplate) => {
    const endpoint = `/contests/${contestId}/problems`;

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

  getCompClasses = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/comp-classes`;

    const result = await this.axiosInstance.get(endpoint);

    return z.array(compClassSchema).parse(result.data);
  };

  createCompClass = async (contestId: number, template: CompClassTemplate) => {
    const endpoint = `/contests/${contestId}/comp-classes`;

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

  getTicks = async (contenderId: number) => {
    const endpoint = `/contenders/${contenderId}/ticks`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(tickSchema).parse(result.data);
  };

  createTick = async (
    contenderId: number,
    tick: Omit<Tick, "id" | "timestamp">,
  ) => {
    const endpoint = `/contenders/${contenderId}/ticks`;

    const result = await this.axiosInstance.post(endpoint, tick, {
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

  getScoreboard = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/scoreboard`;

    const result = await this.axiosInstance.get(endpoint);

    return z.array(scoreboardEntrySchema).parse(result.data);
  };

  getScoreEngines = async (contestId: ContestID) => {
    const endpoint = `/contests/${contestId}/score-engines`;

    const result = await this.axiosInstance.get(endpoint, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return z.array(z.string().uuid()).parse(result.data);
  };

  startScoreEngine = async (
    contestId: ContestID,
    args: StartScoreEngineArguments,
  ) => {
    const endpoint = `/contests/${contestId}/score-engines`;

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
}
