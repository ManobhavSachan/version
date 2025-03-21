"use client"

import { useEffect, useState } from "react"
import { AlertCircle } from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { SystemCard } from "@/app/components/SystemCard"
import { ApplicationsTable } from "@/app/components/ApplicationsTable"
import { CardSkeletons, TableSkeletons } from "@/app/components/LoadingSkeletons"

interface OsVersion {
  name: string
  version: string
  platform: string
}

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

interface ApiResponse {
  os_version: OsVersion
  osquery_version: string
  installed_apps: InstalledApp[]
}

export default function Dashboard() {
  const [data, setData] = useState<ApiResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        const response = await fetch("http://localhost:7070/api/latest_data")

        if (!response.ok) {
          const errorText = await response.text()
          throw new Error(`API request failed with status ${response.status}: ${errorText}`)
        }

        const result = await response.json()
        if (!result.installed_apps || !Array.isArray(result.installed_apps)) {
          throw new Error('Invalid data format received from API')
        }
        setData(result)
      } catch (err) {
        if (err instanceof TypeError && err.message.includes('Failed to fetch')) {
          setError('Cannot connect to API server. Please ensure the server is running at http://localhost:7070')
        } else {
          setError(err instanceof Error ? err.message : "An unknown error occurred")
        }
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  const formatDate = (timestamp: number) => {
    if (!timestamp) return "N/A"
    return new Date(timestamp * 1000).toLocaleString()
  }

  if (error) {
    return (
      <div className="container mx-auto py-10">
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>Failed to fetch system data: {error}</AlertDescription>
        </Alert>
      </div>
    )
  }

  return (
    <div className="container mx-auto py-6 space-y-6">
      <div className="flex flex-col space-y-2">
        <h1 className="text-3xl font-bold tracking-tight">System Information Dashboard</h1>
        <p className="text-muted-foreground">Displaying data collected from osquery on your local system</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {loading ? (
          <CardSkeletons />
        ) : (
          <>
            <SystemCard
              title="Operating System"
              icon="os"
              mainValue={data?.os_version.name || "Unknown"}
              details={[
                { label: "Version", value: data?.os_version.version || "Unknown" },
                { label: "Platform", value: data?.os_version.platform || "Unknown" }
              ]}
            />
            <SystemCard
              title="Osquery Version"
              icon="osquery"
              mainValue={data?.osquery_version || "Unknown"}
            />
            <SystemCard
              title="Applications"
              icon="apps"
              mainValue={data?.installed_apps.length || 0}
            />
          </>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Installed Applications</CardTitle>
          <CardDescription>Applications installed on this system as reported by osquery</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <TableSkeletons />
          ) : (
            <ApplicationsTable
              apps={data?.installed_apps || []}
              formatDate={formatDate}
            />
          )}
        </CardContent>
      </Card>
    </div>
  )
}

