import type { DescMessage } from "@bufbuild/protobuf";

type DirtyMeta = { isDirty?: boolean } | undefined;

// localName (camelCase TS name) -> name (snake_case proto name = mask path).
const cache = new WeakMap<DescMessage, Map<string, string>>();

function protoNameByLocal(schema: DescMessage): Map<string, string> {
  let map = cache.get(schema);
  if (!map) {
    map = new Map(schema.fields.map((f) => [f.localName, f.name]));
    cache.set(schema, map);
  }
  return map;
}

/**
 * Builds protobuf field-mask paths from a form's dirty fields.
 *
 * Reads the proto field names straight from the message schema, so the paths
 * stay in sync with the .proto with no manual maintenance.
 */
export function dirtyFieldMaskPaths(
  schema: DescMessage,
  fieldMeta: Record<string, DirtyMeta>,
): string[] {
  const names = protoNameByLocal(schema);
  return Object.entries(fieldMeta)
    .filter(([, meta]) => meta?.isDirty)
    .map(([key]) => names.get(key))
    .filter((p): p is string => p !== undefined);
}
