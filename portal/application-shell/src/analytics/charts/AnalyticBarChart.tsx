import React from "react"
import { Alert, Button } from "@mui/material"
import { BarChart } from "@mui/x-charts/BarChart"

type AnalyticBarChartProps = {
    labels: string[]
    values: number[]
    seriesLabel: string
    loading: boolean
    error: boolean
    messages: { empty: string; error: string; retry: string }
    retryHandler: () => void
}

const AnalyticBarChart: React.FC<AnalyticBarChartProps> = ({
    labels,
    values,
    seriesLabel,
    loading,
    error,
    messages,
    retryHandler
}) => {
    if (error) {
        return <Alert
            severity="error"
            action={<Button color="inherit" size="small" onClick={retryHandler}>{messages.retry}</Button>}>
            {messages.error}
        </Alert>
    }

    return <BarChart
        xAxis={[{ scaleType: 'band', data: labels }]}
        series={[{
            data: values,
            label: seriesLabel,
            valueFormatter: (value) => value === null ? "" : value.toFixed(2),
            barLabel: "value"
        }]}
        height={400}
        loading={loading}
        localeText={{ noData: messages.empty }}
    />
}

export default AnalyticBarChart
