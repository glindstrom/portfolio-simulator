import React from "react";
import {
    Line,
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer,
    Area,
    AreaChart,
    CartesianGrid,
    Legend,
} from "recharts";

type ChartPoint = {
    month: number;
    median: number;
    p5: number;
    p95: number;
};

type FinalSummary = {
    percentile: string;
    value: number;
}[];

interface Props {
    paths: number[][];
    finalStats: {
        mean: number;
        median: number;
        min: number;
        max: number;
    };
    successRate: number;
}

function percentile(sortedArr: number[], p: number): number {
    const index = (p / 100) * (sortedArr.length - 1);
    const lower = Math.floor(index);
    const upper = Math.ceil(index);
    if (lower === upper) return sortedArr[lower];
    return sortedArr[lower] + (sortedArr[upper] - sortedArr[lower]) * (index - lower);
}

function calculateChartData(paths: number[][]): ChartPoint[] {
    const months = paths[0]?.length || 0;
    const chartData: ChartPoint[] = [];

    for (let i = 0; i < months; i++) {
        const valuesAtMonth = paths.map((p) => p[i]).sort((a, b) => a - b);
        chartData.push({
            month: i,
            median: percentile(valuesAtMonth, 50),
            p5: percentile(valuesAtMonth, 5),
            p95: percentile(valuesAtMonth, 95),
        });
    }

    return chartData;
}

const SimulationResult: React.FC<Props> = ({ paths, finalStats, successRate }) => {
    const chartData = calculateChartData(paths);
    const finalSummary: FinalSummary = [
        { percentile: "Min", value: finalStats.min },
        { percentile: "Median", value: finalStats.median },
        { percentile: "Mean", value: finalStats.mean },
        { percentile: "Max", value: finalStats.max },
        { percentile: "Success Rate", value: successRate * 100 },
    ];

    return (
        <div className="mt-8 space-y-8">
            <div>
                <h2 className="text-xl font-semibold mb-4">Simulation Results</h2>
                <ResponsiveContainer width="100%" height={400}>
                    <AreaChart data={chartData}>
                        <defs>
                            <linearGradient id="colorP5" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#8884d8" stopOpacity={0.3} />
                                <stop offset="95%" stopColor="#8884d8" stopOpacity={0} />
                            </linearGradient>
                        </defs>
                        <XAxis dataKey="month" />
                        <YAxis />
                        <CartesianGrid strokeDasharray="3 3" />
                        <Tooltip />
                        <Legend />
                        <Line type="monotone" dataKey="median" stroke="#4f46e5" strokeWidth={2} dot={false} />
                        <Area type="monotone" dataKey="p5" stroke="#8884d8" fillOpacity={0.2} fill="url(#colorP5)" />
                        <Area type="monotone" dataKey="p95" stroke="#8884d8" fillOpacity={0} />
                    </AreaChart>
                </ResponsiveContainer>
            </div>

            <div className="overflow-x-auto">
                <table className="min-w-full table-auto border border-gray-300 rounded shadow">
                    <thead className="bg-gray-100">
                    <tr>
                        <th className="text-left px-4 py-2 border-b">Statistic</th>
                        <th className="text-right px-4 py-2 border-b">Value</th>
                    </tr>
                    </thead>
                    <tbody>
                    {finalSummary.map((row) => (
                        <tr key={row.percentile} className="border-t">
                            <td className="px-4 py-2">{row.percentile}</td>
                            <td className="px-4 py-2 text-right">
                                {row.percentile === "Success Rate"
                                    ? `${row.value.toFixed(1)}%`
                                    : row.value.toLocaleString(undefined, {
                                        minimumFractionDigits: 0,
                                        maximumFractionDigits: 0,
                                    })}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default SimulationResult;
