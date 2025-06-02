// src/components/Chart.tsx
import {
    LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid,
} from "recharts";
import type { SimulationResponse } from "../types";

interface ChartProps {
    simData: SimulationResponse;
}

// Define a list of distinct colors for the paths
const lineColors = [
    "#8884d8", "#82ca9d", "#ffc658", "#ff7300", "#00C49F",
    "#FFBB28", "#FF8042", "#0088FE", "#AF19FF", "#FF1919",
];

export default function Chart({ simData }: ChartProps) {
    const pathLimit = 10; // Limit how many paths to show
    const paths = simData.paths.slice(0, pathLimit);

    // Ensure paths are not empty and have data before proceeding
    if (!paths || paths.length === 0 || paths[0].length === 0) {
        return <div className="w-full h-[400px] flex items-center justify-center">No data to display</div>;
    }

    const periods = paths[0].length;

    // Build data format for Recharts
    const chartData = Array.from({ length: periods }).map((_, i) => {
        const point: any = { year: Math.floor(i / 12) }; // whole years
        paths.forEach((path, idx) => {
            point[`Path ${idx + 1}`] = path[i];
        });
        return point;
    });

    // Filter out only 1y, 5y, 10y, ... ticks
    const uniqueYears = [...new Set(chartData.map(d => d.year))];
    // Ensure year 0 is included if present, and year 1, then multiples of 5
    const xTicks = uniqueYears.filter((y) => y === 0 || y === 1 || y % 5 === 0);
    // If year 1 isn't naturally in uniqueYears (e.g. less than 2 years data), add it if relevant.
    // The current logic is fine for longer periods.

    return (
        <div className="w-full h-[400px] bg-white p-4 rounded-lg shadow-md"> {/* Added some basic styling to the container */}
            <ResponsiveContainer width="100%" height="100%">
                <LineChart
                    data={chartData}
                    margin={{
                        top: 5,
                        right: 30, // Increased right margin for potential last tick label
                        left: 35,  // Increased left margin for Y-axis label
                        bottom: 20, // Increased bottom margin for X-axis label
                    }}
                >
                    <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" vertical={false} />
                    <XAxis
                        dataKey="year"
                        ticks={xTicks}
                        tickFormatter={(val) => `${val}y`}
                        tick={{ fontSize: '12px', fill: '#666' }}
                        label={{
                            value: "Year",
                            position: "outsideCenter",
                            dy: 15, // Pushes label further down from axis
                            style: { fill: "#555", fontSize: '14px', fontWeight: 'bold' },
                        }}
                        axisLine={{ stroke: '#ccc' }}
                        tickLine={{ stroke: '#ccc' }}
                    />
                    <YAxis
                        tickFormatter={(val) => `$${(val / 1e6).toFixed(1)}M`}
                        tickCount={7} // Suggest number of ticks to prevent overlap
                        domain={[0, 'auto']} // Start Y-axis from 0
                        tick={{ fontSize: '12px', fill: '#666' }}
                        label={{
                            value: "Portfolio Value",
                            angle: -90,
                            position: "left", // Position outside, to the left of ticks
                            offset: 20,     // Distance from the Y-axis line/ticks
                            style: { textAnchor: 'middle', fill: "#555", fontSize: '14px', fontWeight: 'bold' },
                        }}
                        axisLine={{ stroke: '#ccc' }}
                        tickLine={{ stroke: '#ccc' }}
                    />
                    <Tooltip
                        formatter={(value: number, name: string) => [`$${value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`, name]}
                        labelFormatter={(label) => `Year: ${label}`}
                        cursor={{ stroke: '#8884d8', strokeWidth: 1, strokeDasharray: '3 3' }}
                        contentStyle={{
                            backgroundColor: 'rgba(255, 255, 255, 0.9)',
                            border: '1px solid #ccc',
                            borderRadius: '8px',
                            boxShadow: '0px 2px 10px rgba(0,0,0,0.1)',
                            padding: '10px'
                        }}
                        labelStyle={{ fontWeight: 'bold', color: '#333', marginBottom: '5px' }}
                        itemStyle={{ color: '#555' }}
                    />
                    {paths.map((_, idx) => (
                        <Line
                            key={idx}
                            type="monotone"
                            dataKey={`Path ${idx + 1}`}
                            stroke={lineColors[idx % lineColors.length]} // Use distinct colors
                            dot={false}
                            strokeWidth={1.5} // Slightly thicker lines
                            activeDot={{ r: 5, strokeWidth: 0 }} // Style for dot on hover
                        />
                    ))}
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
}