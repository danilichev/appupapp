import { Middleware } from "openapi-fetch";

import { AuthTokens, AuthTokenService } from "./types";

type CreateAuthMiddlewareProps = {
  authTokenService: AuthTokenService;
  baseUrl: string;
  getAuthTokenPaths: string[];
  notProtectedPaths: string[];
  refreshAuthTokenPath: string;
};

const refreshAuthToken = async (
  refreshTokenUrl: string,
  refreshToken?: string,
): Promise<AuthTokens> => {
  if (!refreshToken) throw new Error("No refresh token provided");

  const response = await fetch(refreshTokenUrl, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${refreshToken}`,
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) throw new Error("Refresh auth tokens failed");

  const data = (await response.json()) as AuthTokens;

  if (!data.accessToken || !data.refreshToken) {
    throw new Error("Invalid token response");
  }

  return data;
};

let refreshPromise: Promise<void> | null = null;

export const createAuthMiddleware = ({
  authTokenService,
  baseUrl,
  getAuthTokenPaths,
  notProtectedPaths,
  refreshAuthTokenPath,
}: CreateAuthMiddlewareProps): Middleware => ({
  async onRequest({ request, schemaPath }) {
    if (notProtectedPaths.includes(schemaPath)) return;

    if (refreshPromise) {
      await refreshPromise;
    }

    const { accessToken } = (await authTokenService.getTokens()) || {};

    if (!accessToken) {
      throw new Error("No access token provided");
    }

    request.headers.set("Authorization", `Bearer ${accessToken}`);

    return request;
  },
  async onResponse({ response, request, schemaPath }) {
    if (response.ok) {
      if (getAuthTokenPaths.includes(schemaPath)) {
        const data = (await response.clone().json()) as AuthTokens;

        if (!data.accessToken || !data.refreshToken) {
          throw new Error("Invalid token response");
        }

        authTokenService.setTokens(data);
      }

      return response;
    }

    if (response.status === 401) {
      if (refreshPromise) {
        await refreshPromise;
      } else {
        refreshPromise = (async () => {
          try {
            const { refreshToken } = (await authTokenService.getTokens()) || {};
            const refreshUrl = `${baseUrl}${refreshAuthTokenPath}`;
            const newTokens = await refreshAuthToken(refreshUrl, refreshToken);

            authTokenService.setTokens(newTokens);
          } catch {
            authTokenService.resetTokens();
            throw new Error("Failed to refresh auth tokens");
          } finally {
            refreshPromise = null;
          }
        })();

        await refreshPromise;
      }

      const { accessToken } = (await authTokenService.getTokens()) || {};

      request.headers.set("Authorization", `Bearer ${accessToken}`);

      const res = await fetch(request.url, {
        body: request.body,
        headers: request.headers,
        method: request.method,
      });

      return res;
    }
  },
});
