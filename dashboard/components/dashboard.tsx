"use client";

import { LineChart } from "@mui/x-charts";
import { useSession } from "next-auth/react";
import { useEffect, useState } from "react";

export function Dashboard() {
    const { data } = useSession()
    const [id, token] = [data?.id!, data?.token!]

    const [currentDashboard, setCurrentDashboard] = useState("test")

    return (
        <div className="text-gray-100">
            <div>
                <p>TestDashboard</p>
            </div>
            <div className="h-[1px] bg-gray-300 w-full"></div>
            <TestDashboard id={id} token={token} />
        </div>
    )
}

interface GroupedRecord {
    timestamp: string;
    window: number;
    event: string;
    count: number
}

function TestDashboard({ id, token }: { id: string, token: string }) {
    const [yAxisData, setYAxisData] = useState<number[]>([]) // timestamp vs count
    const [xAxisData, setXAxisData] = useState<Date[]>([]) 

    useEffect(() => {
        const fetchData = async () => {
            if(!id || !token) return;
            const resp = await fetch("http://localhost:8080/analytics?event=testEvent&window=3600", {
                headers: {
                    "Authorization": id + ":" + token
                }
            })

            if(resp.ok) {
                let xAxis: Date[] = []
                let yAxis: number[] = []
                const body: GroupedRecord[] = await resp.json()
                body.forEach(item => {
                    xAxis.push(new Date(item.timestamp))
                    yAxis.push(item.count)
                })

                setXAxisData(xAxis)
                setYAxisData(yAxis)
            }
        }

        fetchData()
    }, [id, token])

    return (
        <div>
            <LineChart
                className="border-white border-2 bg-white"

                xAxis={
                    [{ data: xAxisData, scaleType: "time" }]
                }

                series={[
                    {
                        data: yAxisData,
                        showMark: false
                    },
                ]}

                height={300}
                grid={{ vertical: true, horizontal: true }}
            />
        </div>
    )
}