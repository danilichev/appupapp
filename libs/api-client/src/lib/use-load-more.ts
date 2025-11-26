import { UseInfiniteQueryResult } from "@tanstack/react-query";
import { useCallback } from "react";

type UseLoadMoreProps = Pick<
  UseInfiniteQueryResult,
  "hasNextPage" | "isFetching"
> & {
  fetchNextPage: UseInfiniteQueryResult["fetchNextPage"] | (() => void);
};

export const useLoadMore = ({
  fetchNextPage,
  hasNextPage,
  isFetching,
}: UseLoadMoreProps) => {
  const loadMore = useCallback(() => {
    if (hasNextPage && !isFetching) {
      fetchNextPage();
    }
  }, [fetchNextPage, hasNextPage, isFetching]);

  return loadMore;
};
