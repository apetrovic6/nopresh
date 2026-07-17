import { getBloodPressure } from "#/gen/proto/bloodpressure/v1/bloodpressure-BloodPressureService_connectquery"
import { useInfiniteQuery } from "@connectrpc/connect-query"
import { Spinner } from "../ui/spinner";
import { BloodPressureItem } from "./blood-pressure-item";
import { SORTORDER, type BloodPressureEntry } from "#/gen/proto/bloodpressure/v1/bloodpressure_pb";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { format, parseISO } from "date-fns";
import { useVirtualizer } from "@tanstack/react-virtual";
import { useEffect, useMemo, useRef, useState } from "react";
import { useAppForm } from "#/hooks/demo.form";
import type { DateRange } from "react-day-picker";
import { Button } from "../ui/button";

const PAGE_SIZE = 20;

// A flat, virtualizable row: either a date header or a single entry.
type Row =
  | { kind: "header"; date: string; key: string }
  | { kind: "entry"; entry: BloodPressureEntry; key: string };

interface BloodPressureTableProps {
  onEdit: (entry: BloodPressureEntry) => void;
  onDelete: (entry: BloodPressureEntry) => void;
}

export function BloodPressureTable({ onEdit, onDelete }: BloodPressureTableProps) {
  // appliedRange is the committed filter the query actually uses. The form
  // holds the range. It only gets set on Apply (onSubmit), so
  // picking dates doesn't refetch until the user applies the filter.
  const [appliedRange, setAppliedRange] = useState<DateRange | undefined>(undefined);

  const filterForm = useAppForm({
    defaultValues: {
      dateRange: undefined as DateRange | undefined,
    },
    onSubmit: ({ value }) => {
      // This changes the query key, so the connect-query refetches automatically
      setAppliedRange(value.dateRange);
    },
  });

  const parentRef = useRef<HTMLDivElement>(null);

  // Step 1: infinite query. `pageInfo` is the request's page param; its `cursor`
  // starts empty and is replaced per page by getNextPageParam.
  // Note: connect-query's useInfiniteQuery locks `select` to InfiniteData, so we
  // transform the data in a useMemo below instead of via `select`.
  const { data, isLoading, hasNextPage, isFetchingNextPage, fetchNextPage } = useInfiniteQuery(
    getBloodPressure,
    {
      dateFilter: {
        startDate: appliedRange?.from?.toTimestamp(),
        endDate: appliedRange?.to?.toTimestamp(),
      },
      pageInfo: { cursor: "", pageSize: PAGE_SIZE, sortOrder: SORTORDER.SORTORDER_DESC }
    },
    {
      pageParamKey: "pageInfo",
      getNextPageParam: (lastPage) =>
        lastPage.pageInfo?.cursor
          ? { cursor: lastPage.pageInfo.cursor, pageSize: PAGE_SIZE, sortOrder: SORTORDER.SORTORDER_DESC }
          : undefined,
    },
  );

  // Step 2: flatten InfiniteData -> grouped -> a single list of header/entry rows.
  // Backend already returns newest-first (date desc, id desc); groupBy preserves
  // that insertion order for both the day groups and the entries within them.
  const rows = useMemo<Row[]>(() => {
    const entries = data?.pages.flatMap((page) => page.bloodPressure) ?? [];

    const groups = Map.groupBy(entries, (item) =>
      format(timestampDate(item.dateTimeUtc!), "yyyy-MM-dd"),
    );

    const result: Row[] = [];
    for (const [date, items] of groups) {
      result.push({ kind: "header", date, key: `header-${date}` });
      for (const entry of items) {
        result.push({ kind: "entry", entry, key: `entry-${entry.id}` });
      }
    }
    return result;
  }, [data]);

  // Step 3: virtualize the flat rows. estimateSize is a rough guess; real heights
  // are measured from the DOM via measureElement below.
  const virtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 140,
    overscan: 8,
    getItemKey: (index) => rows[index].key,
  });

  const virtualItems = virtualizer.getVirtualItems();

  // Step 5: load the next page once the last rows enter (or near) the viewport.
  useEffect(() => {
    const last = virtualItems[virtualItems.length - 1];
    if (!last) return;

    if (last.index >= rows.length - 1 && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [virtualItems, rows.length, hasNextPage, isFetchingNextPage, fetchNextPage]);

  function onSubmitFilters(e: React.SubmitEvent) {
    e.preventDefault();
    e.stopPropagation();

    filterForm.handleSubmit();
  }


  if (isLoading) {
    return <Spinner />;
  }

  // Step 4: a bounded, scrollable container; a spacer of total height; each
  // virtual row absolutely positioned and self-measured.
  return (
    <>
      <form className="flex justify-between my-2 px-4" onSubmit={onSubmitFilters}>
        <filterForm.AppField name="dateRange">
          {field => <field.DateRangePicker label="Date Range" onClear={() => setAppliedRange(undefined)} fieldGroupClass="max-w-xs" />}
        </filterForm.AppField>

        <filterForm.Subscribe selector={state => state.values.dateRange}>
          {(state) =>
            <Button disabled={state?.from === undefined || state?.to === undefined} className="self-end" type="submit">
              Apply
            </Button>
          }
        </filterForm.Subscribe>
      </form >

      <div ref={parentRef} className="flex-1 min-h-0 overflow-auto px-4">
        <div style={{ height: virtualizer.getTotalSize(), position: "relative" }}>
          {virtualItems.map((virtualItem) => {
            const row = rows[virtualItem.index];
            return (
              <div
                key={virtualItem.key}
                data-index={virtualItem.index}
                ref={virtualizer.measureElement}
                className="pb-2"
                style={{
                  position: "absolute",
                  top: 0,
                  left: 0,
                  width: "100%",
                  transform: `translateY(${virtualItem.start}px)`,
                }}
              >
                {row.kind === "header" ? (
                  <h2 className="text-sm font-semibold text-muted-foreground my-2 sticky">
                    {format(parseISO(row.date), "PPP")}
                  </h2>
                ) : (
                  <BloodPressureItem item={row.entry} onEdit={onEdit} onDelete={onDelete} />
                )}
              </div>
            );
          })}
        </div>

        {isFetchingNextPage && (
          <div className="flex justify-center py-4">
            <Spinner />
          </div>
        )}
      </div>
    </>
  );
}
