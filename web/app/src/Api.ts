import type { RawAxiosRequestHeaders } from "axios";
import axios from "axios";
import configData from "./config.json";
import type { Contender } from "./models/contender";
import type { Contest } from "./models/contest";
import type { Problem } from "./models/problem";
import type { CompClass } from "./models/compClass";
import type { Tick } from "./models/tick";
import type { Scoreboard } from "./models/scoreboard";

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
  private static baseUrl: string = configData.API_URL;
  private credentialsProvider: ApiCredentialsProvider | undefined;

  private constructor() { }

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
    const endpoint = `/contenders/findByCode?code=${registrationCode}`;

    const result = await axios.get<Contender>(
      `${ApiClient.baseUrl}${endpoint}`,
      {
        headers: this.credentialsProvider?.getAuthHeaders(),
      }
    );

    return result.data;
  };

  getContender = async (id: number) => {
    const endpoint = `/contenders/${id}`;

    const result = await axios.get<Contender>(
      `${ApiClient.baseUrl}${endpoint}`,
      {
        headers: this.credentialsProvider?.getAuthHeaders(),
      }
    );

    return result.data;
  };

  updateContender = async (id: number, contender: Contender) => {
    const endpoint = `/contenders/${id}`;

    const result = await axios.put<Contender>(
      `${ApiClient.baseUrl}${endpoint}`,
      contender,
      {
        headers: this.credentialsProvider?.getAuthHeaders(),
      }
    );

    return result.data;
  };

  getContest = async (id: number) => {
    const endpoint = `/contests/${id}`;

    const result = await axios.get<Contest>(`${ApiClient.baseUrl}${endpoint}`, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return result.data;
  };

  getProblems = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/problems`;

    const result = await axios.get<Problem[]>(
      `${ApiClient.baseUrl}${endpoint}`,
      {
        headers: this.credentialsProvider?.getAuthHeaders(),
      }
    );

    return result.data;
  };

  getCompClasses = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/compClasses`;

    const result = await axios.get<CompClass[]>(
      `${ApiClient.baseUrl}${endpoint}`,
      {
        headers: this.credentialsProvider?.getAuthHeaders(),
      }
    );

    return result.data;
  };

  getTicks = async (contenderId: number) => {
    const endpoint = `/contenders/${contenderId}/ticks`;

    const result = await axios.get<Tick[]>(`${ApiClient.baseUrl}${endpoint}`, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });

    return result.data;
  };

  createTick = async (contenderId: number, tick: Omit<Tick, "id">) => {
    const endpoint = `/contenders/${contenderId}/ticks`;

    const result = await axios.post<Tick>(
      `${ApiClient.baseUrl}${endpoint}`,
      tick,
      {
        headers: this.credentialsProvider?.getAuthHeaders(),
      }
    );

    return result.data;
  };

  deleteTick = async (tickId: number) => {
    const endpoint = `/ticks/${tickId}`;

    await axios.delete(`${ApiClient.baseUrl}${endpoint}`, {
      headers: this.credentialsProvider?.getAuthHeaders(),
    });
  };

  getScoreboard = async (contestId: number) => {
    const endpoint = `/contests/${contestId}/scoreboard`;

    const result = await axios.get<Scoreboard>(
      `${ApiClient.baseUrl}${endpoint}`
    );

    return result.data;
  };
}
