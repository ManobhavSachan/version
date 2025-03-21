import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { CheckCircle2, Database, Laptop, Package } from "lucide-react"

interface SystemCardProps {
  title: string
  icon: "os" | "osquery" | "apps"
  mainValue: string | number
  details?: {
    label: string
    value: string
  }[]
}

export function SystemCard({ title, icon, mainValue, details }: SystemCardProps) {
  const icons = {
    os: <Laptop className="h-4 w-4 text-muted-foreground" />,
    osquery: <Database className="h-4 w-4 text-muted-foreground" />,
    apps: <Package className="h-4 w-4 text-muted-foreground" />
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icons[icon]}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{mainValue}</div>
        {icon === "osquery" && (
          <div className="flex items-center mt-4">
            <CheckCircle2 className="mr-2 h-4 w-4 text-green-500" />
            <span className="text-sm text-muted-foreground">Active and running</span>
          </div>
        )}
        {icon === "apps" && (
          <p className="text-xs text-muted-foreground mt-4">
            Total installed applications detected on this system
          </p>
        )}
        {details?.map((detail, index) => (
          <div key={index} className="flex items-center justify-between mt-2">
            <CardDescription>{detail.label}</CardDescription>
            <Badge variant="outline">{detail.value}</Badge>
          </div>
        ))}
      </CardContent>
    </Card>
  )
} 