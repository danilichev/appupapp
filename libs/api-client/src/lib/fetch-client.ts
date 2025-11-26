import createOpenApiFetchClient from "openapi-fetch";

import { paths } from "./api";
import { createAuthMiddleware } from "./auth-middleware";
import { AuthTokenService } from "./types";

export type CreateFetchClientProps = {
  authTokenService: AuthTokenService;
  baseUrl: string;
};

const getAuthTokenPaths = ["/auth/login", "/auth/register"];

export type FetchClient = ReturnType<typeof createOpenApiFetchClient<paths>>;

export const createFetchClient = ({
  authTokenService,
  baseUrl,
}: CreateFetchClientProps): FetchClient => {
  const fetchClient = createOpenApiFetchClient<paths>({ baseUrl });

  const authMiddleware = createAuthMiddleware({
    authTokenService,
    baseUrl,
    getAuthTokenPaths,
    notProtectedPaths: getAuthTokenPaths,
    refreshAuthTokenPath: "/auth/refresh",
  });

  fetchClient.use(authMiddleware);

  return fetchClient;
};
