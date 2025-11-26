import createReactQueryClient from "openapi-react-query";

import { paths } from "./api";
import { FetchClient } from "./fetch-client";

export type ApiClient = ReturnType<typeof createReactQueryClient<paths>>;

export const createApiClient = (fetchClient: FetchClient) =>
  createReactQueryClient(fetchClient);
