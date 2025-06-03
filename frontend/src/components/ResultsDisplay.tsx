// src/components/ResultsDisplay.tsx
import React from "react";
import {
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer,
    Area,
    AreaChart,
    CartesianGrid,
    Legend,
    ReferenceLine,
} from "recharts";
import type { SimulationResponse as ResultPayload } from "../types";
import {formatUSD, yAxisTickFormatter, tooltipFormatterFunc, xAxisTickFormatter} from "../utils/formatters";
import {calculateChartData, type ChartPoint, generateXAxisYearTicks} from "../utils/chartUtils"; // Assuming percentile is not directly used here anymore

// SummaryRow type definition remains local as it's specific to this component's table structure.
type SummaryRow = {
    label: string;
    displayValue: string;
};

interface Props {
    resultData: ResultPayload;
}

const ResultsDisplay: React.FC<Props> = ({ resultData }) => {
    const { paths, finalStats, successRate, simulatedCAGR } = resultData;

    const chartData: ChartPoint[] = calculateChartData(paths);
    const finalSummary: SummaryRow[] = [
        { label: "Min Final Value", displayValue: formatUSD(finalStats.min) },
        { label: "Median Final Value", displayValue: formatUSD(finalStats.median) },
        { label: "Mean Final Value", displayValue: formatUSD(finalStats.mean) },
        { label: "Max Final Value", displayValue: formatUSD(finalStats.max) },
    ];
    if (typeof simulatedCAGR === 'number' && !isNaN(simulatedCAGR)) {
        finalSummary.push({ label: "Simulated CAGR", displayValue: `${(simulatedCAGR * 100).toFixed(2)}%` });
    }
    finalSummary.push({ label: "Success Rate", displayValue: `${(successRate * 100).toFixed(1)}%` });

    // Use the imported utility function to generate ticks
    const xTicks: number[] = generateXAxisYearTicks(chartData.length);

    return (
        <div className="bg-white shadow-xl rounded-xl p-6 sm:p-8 space-y-6">
            <h2 className="text-xl sm:text-2xl font-semibold text-gray-700 mb-4">Simulation Results</h2>
            <div>
                <h3 className="text-lg sm:text-xl font-semibold mb-3 sm:mb-4 text-gray-700">Projection Distribution</h3>
                <ResponsiveContainer width="100%" height={400}>
                    <AreaChart data={chartData} margin={{ top: 5, right: 30, left: 20, bottom: 25 }}>
                        <defs>
                            <linearGradient id="gradientMedianResults" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#4f46e5" stopOpacity={0.15}/>
                                <stop offset="95%" stopColor="#4f46e5" stopOpacity={0}/>
                            </linearGradient>
                            <linearGradient id="gradientPercentileRangeResults" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#8884d8" stopOpacity={0.2}/>
                                <stop offset="95%" stopColor="#8884d8" stopOpacity={0.05}/>
                            </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                        <XAxis
                            dataKey="month"
                            type="number"
                            domain={[0, 'dataMax']}
                            ticks={xTicks.length > 0 ? xTicks : undefined}
                            tickFormatter={xAxisTickFormatter} // Use imported formatter
                            label={{ value: "Year", position: 'insideBottom', offset: 0, dy: 10, fontSize: 12, fill: '#555' }}
                            tick={{ fontSize: 11, fill: '#666' }}
                        />
                        <YAxis
                            tickFormatter={yAxisTickFormatter}
                            width={80}
                            tick={{ fontSize: 11, fill: '#666' }}
                            label={{ value: 'Portfolio Value', angle: -90, position: 'insideLeft', dx: -15, fontSize: 12, fill: '#555' }}
                        />
                        <Tooltip formatter={tooltipFormatterFunc} />
                        <Legend wrapperStyle={{paddingTop: '15px'}}/>
                        <Area type="monotone" dataKey="p5_p95_range" fill="url(#gradientPercentileRangeResults)" stroke="none" activeDot={false} name="5th - 95th Percentile" />
                        <Area type="monotone" dataKey="median" stroke="#4f46e5" fill="url(#gradientMedianResults)" strokeWidth={2.5} dot={false} name="Median Value" />
                        <ReferenceLine y={0} stroke="#666" strokeWidth={1} strokeDasharray="2 2" ifOverflow="extendDomain"/>
                    </AreaChart>
                </ResponsiveContainer>
            </div>

            {/* Table Section */}
            <div className="overflow-x-auto mt-6">
                <h2 className="text-lg sm:text-xl font-semibold mb-3 sm:mb-4 text-gray-700">Summary Statistics</h2>
                <table className="min-w-full w-full table-auto border border-gray-200 rounded-lg shadow-sm bg-white">
                    <thead className="bg-gray-50">
                    <tr>
                        <th className="text-left px-4 py-3 border-b border-gray-200 text-xs font-semibold text-gray-600 uppercase tracking-wider">Statistic</th>
                        <th className="text-right px-4 py-3 border-b border-gray-200 text-xs font-semibold text-gray-600 uppercase tracking-wider">Value</th>
                    </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                    {finalSummary.map((row) => (
                        <tr key={row.label} className="hover:bg-gray-50">
                            <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-700">{row.label}</td>
                            <td className="px-4 py-3 whitespace-nowrap text-right text-sm text-gray-800 font-medium">
                                {row.displayValue}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default ResultsDisplay;