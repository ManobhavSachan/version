import { Package } from "lucide-react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

interface InstalledApp {
  name: string
  path: string
  bundle_identifier: string
  bundle_name: string
  bundle_short_version: string
  display_name: string
  minimum_system_version?: string
  last_opened_time: number
  end_time?: number
}

interface ApplicationsTableProps {
  apps: InstalledApp[]
  formatDate: (timestamp: number) => string
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="flex flex-col items-center justify-center p-8 text-center">
      <Package className="h-8 w-8 text-muted-foreground mb-4" />
      <p className="text-lg font-medium">No Applications Found</p>
      <p className="text-sm text-muted-foreground">{message}</p>
    </div>
  )
}

export function ApplicationsTable({ apps, formatDate }: ApplicationsTableProps) {
  const recentApps = apps
    .sort((a, b) => (b.last_opened_time || 0) - (a.last_opened_time || 0))
    .slice(0, 10)

  const deletedApps = apps
    .filter(app => app.end_time != null)
    .sort((a, b) => (b.end_time || 0) - (a.end_time || 0))

  return (
    <Tabs defaultValue="all" className="w-full">
      <TabsList className="mb-4">
        <TabsTrigger value="all">All Applications</TabsTrigger>
        <TabsTrigger value="recent">Recently Used</TabsTrigger>
        <TabsTrigger value="deleted">Deleted Apps</TabsTrigger>
      </TabsList>

      <TabsContent value="all" className="w-full">
        <div className="rounded-md border">
          {apps.length === 0 ? (
            <EmptyState message="There are no applications installed on this system." />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Bundle ID</TableHead>
                  <TableHead>Version</TableHead>
                  <TableHead className="hidden md:table-cell">Path</TableHead>
                  <TableHead className="hidden lg:table-cell">Last Opened</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {apps.map((app, index) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium">
                      {app.display_name || app.bundle_name || app.name}
                    </TableCell>
                    <TableCell>{app.bundle_identifier}</TableCell>
                    <TableCell>{app.bundle_short_version || "N/A"}</TableCell>
                    <TableCell className="hidden md:table-cell truncate max-w-[200px]" title={app.path}>
                      {app.path}
                    </TableCell>
                    <TableCell className="hidden lg:table-cell">{formatDate(app.last_opened_time)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
      </TabsContent>

      <TabsContent value="recent">
        <div className="rounded-md border">
          {!apps.some(app => app.last_opened_time) ? (
            <EmptyState message="No applications have been opened recently." />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Bundle ID</TableHead>
                  <TableHead>Version</TableHead>
                  <TableHead className="hidden md:table-cell">Last Opened</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recentApps.map((app, index) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium">
                      {app.display_name || app.bundle_name || app.name}
                    </TableCell>
                    <TableCell>{app.bundle_identifier}</TableCell>
                    <TableCell>{app.bundle_short_version || "N/A"}</TableCell>
                    <TableCell className="hidden md:table-cell">{formatDate(app.last_opened_time)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
      </TabsContent>

      <TabsContent value="deleted">
        <div className="rounded-md border">
          {deletedApps.length === 0 ? (
            <EmptyState message="There are no deleted applications to display." />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Bundle ID</TableHead>
                  <TableHead>Version</TableHead>
                  <TableHead className="hidden md:table-cell">Last Opened</TableHead>
                  <TableHead className="hidden md:table-cell">Deleted At</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {deletedApps.map((app, index) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium">
                      {app.display_name || app.bundle_name || app.name}
                    </TableCell>
                    <TableCell>{app.bundle_identifier}</TableCell>
                    <TableCell>{app.bundle_short_version || "N/A"}</TableCell>
                    <TableCell className="hidden md:table-cell">{formatDate(app.last_opened_time)}</TableCell>
                    <TableCell className="hidden md:table-cell">{formatDate(app.end_time || 0)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
      </TabsContent>
    </Tabs>
  )
} 