import {
  QueryClient,
  QueryClientConfig,
  QueryClientProvider as Provider,
} from "@tanstack/react-query";
import { PropsWithChildren } from "react";

export const defaultQueryClientConfig: QueryClientConfig = {
  defaultOptions: {
    queries: {
      gcTime: 1000 * 60 * 60 * 24,
      retry: 3,
      staleTime: Infinity,
    },
    mutations: {
      retry: 2,
    },
  },
};

export const defaultQueryClient = new QueryClient(defaultQueryClientConfig);

export type QueryClientProviderProps = PropsWithChildren<{
  client?: QueryClient;
}>;

export const QueryClientProvider = ({
  children,
  client = defaultQueryClient,
}: QueryClientProviderProps) => {
  return <Provider client={client}>{children}</Provider>;
};
