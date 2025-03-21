import { Skeleton } from "@/components/ui/skeleton"

export function CardSkeletons() {
  return (
    <>
      <Skeleton className="h-[180px] w-full" />
      <Skeleton className="h-[180px] w-full" />
      <Skeleton className="h-[180px] w-full" />
    </>
  )
}

export function TableSkeletons() {
  return (
    <div className="space-y-2">
      <Skeleton className="h-4 w-full" />
      <Skeleton className="h-4 w-full" />
      <Skeleton className="h-4 w-full" />
      <Skeleton className="h-4 w-full" />
      <Skeleton className="h-4 w-full" />
    </div>
  )
} 