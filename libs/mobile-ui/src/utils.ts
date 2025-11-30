export const keyExtractor = (
  item: { id?: unknown; key?: unknown },
  index: number,
) => `${item.id ?? item.key ?? index}`;
