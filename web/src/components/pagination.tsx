import { Button } from "./ui/button";
import { PaginationEllipsis, PaginationItem, PaginationLink } from "./ui/pagination";

function getVisiblePages(current: number, total: number): (number | 'ellipsis')[] {
  const pages: (number | 'ellipsis')[] = [];

  const start = Math.max(1, current - 1);
  const end = Math.min(total, current + 1);

  if (start > 1) pages.push('ellipsis');

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }

  if (end < total) pages.push('ellipsis');

  return pages
}

interface PaginationProps {
  currentPage: number,
  pageCount: number,
  setPage: (page: number) => void,
}

export function PaginationPages({ currentPage, pageCount, setPage }: PaginationProps) {
  return getVisiblePages(currentPage, pageCount).map((page, i) =>
    page === 'ellipsis'
      ? <PaginationItem key={`ellipsis-${i}`}><PaginationEllipsis /></PaginationItem>
      : <PaginationItem key={page}>
          <Button variant="ghost" asChild onClick={() => setPage(page - 1)}>
            <PaginationLink isActive={page === currentPage}>{page}</PaginationLink>
          </Button>
        </PaginationItem>
  );
}

