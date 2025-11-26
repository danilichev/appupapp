import { FetchOptions, FetchResponse } from "openapi-fetch";

import { components } from "./api";

export type ApiResponse<T extends Record<string | number, unknown>> =
  FetchResponse<T, FetchOptions<T>, `${string}/${string}`>;

export type AuthTokens = components["schemas"]["AuthToken"];

type AuthTokenResponse = AuthTokens | null;

export type AuthTokenService = {
  getTokens: () => AuthTokenResponse | Promise<AuthTokenResponse>;
  resetTokens: () => void | Promise<void>;
  setTokens: (authTokens: AuthTokens) => void | Promise<void>;
};

export type Folder = components["schemas"]["Folder"];
export type Link = components["schemas"]["Link"];
export type User = components["schemas"]["User"];
