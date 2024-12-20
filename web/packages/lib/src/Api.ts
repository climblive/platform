import type { AxiosInstance, RawAxiosRequestHeaders } from "axios";
import axios from "axios";
import { z } from "zod";
import { contestSchema, scoreboardEntrySchema } from "./models";
import { compClassSchema } from "./models/compClass";
import { contenderSchema, type Contender } from "./models/contender";
import { problemSchema } from "./models/problem";
import { tickSchema, type Tick } from "./models/tick";
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

  updateContender = async (id: number, contender: Contender) => {
    const endpoint = `/contenders/${id}`;

    const result = await this.axiosInstance.put(endpoint, contender, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return contenderSchema.parse(result.data);
  };

  getContest = async (id: number) => {
    const endpoint = `/contests/${id}`;

    const result = await this.axiosInstance.get(endpoint);

    return contestSchema.parse(result.data);
  };

  getProblems = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/problems`;

    const result = await this.axiosInstance.get(endpoint);

    return z.array(problemSchema).parse(result.data);
  };

  getCompClasses = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/compClasses`;

    const result = await this.axiosInstance.get(endpoint);

    return z.array(compClassSchema).parse(result.data);
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
}
