// src/components/Result.tsx
import type { SimulationResponse as Result } from "../types";
import Chart from "./Chart";

type Props = {
    result: Result;
};

export default function ResultView({ result }: Props) {
    return (
        <div className="space-y-4 mt-4">
            {/* Pass the full result to Chart */}
            <Chart simData={result} />

            <div className="bg-gray-100 p-4 rounded text-sm">
                <p><strong>Mean:</strong> ${formatMillions(result.final_stats.mean)}</p>
                <p><strong>Median:</strong> ${formatMillions(result.final_stats.median)}</p>
                <p><strong>Min:</strong> ${formatMillions(result.final_stats.min)}</p>
                <p><strong>Max:</strong> ${formatMillions(result.final_stats.max)}</p>
                <p><strong>Success rate:</strong> {(result.success_rate * 100).toFixed(1)}%</p>
            </div>
        </div>
    );
}

// Format large numbers as e.g. "12.4M"
function formatMillions(value: number): string {
    return `${(value / 1e6).toFixed(2)}M`;
}
