// src/App.tsx
import { useState } from "react";
import Form from "./components/Form";
import Chart from "./components/Chart";
import type { SimulationResponse, Params, SummaryStats } from "./types";

export default function App() {
    const [result, setResult] = useState<SimulationResponse | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const handleSimulate = async (params: Params) => {
        setLoading(true);
        setError(null);
        setResult(null);
        console.log("Sending simulation request:", params);
        try {
            const apiUrl = "http://localhost:8085/api/simulate"; // TODO: Use environment variable

            const res = await fetch(apiUrl, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(params),
            });

            if (!res.ok) {
                const errorData = await res.json().catch(() => ({ message: "An unknown error occurred" }));
                const errorMessage = errorData.detail || errorData.message || `HTTP error! status: ${res.status}`;
                throw new Error(errorMessage);
            }

            const data: SimulationResponse = await res.json();
            console.log("Received result:", data);
            setResult(data);
        } catch (err: any) {
            console.error("Simulation request failed:", err);
            setError(err.message || "Failed to fetch simulation results.");
        } finally {
            setLoading(false);
        }
    };

    const formatCurrency = (value: number) => {
        return value.toLocaleString(undefined, {
            style: "currency",
            currency: "EUR",
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        });
    };

    return (
        <div className="min-h-screen bg-gray-100 p-4 sm:p-6 font-sans">
            {/* Increased max-width for better two-column layout */}
            <div className="max-w-7xl mx-auto space-y-8">
                <header className="text-center py-6">
                    <h1 className="text-3xl sm:text-4xl font-bold text-gray-800">
                        📈 Portfolio Growth Simulator
                    </h1>
                </header>

                {/* Flex container for two-column layout on large screens */}
                <div className="lg:flex lg:gap-8 space-y-8 lg:space-y-0">

                    {/* Column 1: Parameters Form */}
                    {/* On large screens, this div takes 1/3 of the width. On smaller screens, it's full width. */}
                    <div className="lg:w-2/5 xl:w-1/3"> {/* Adjusted width for form: 2/5 on lg, 1/3 on xl */}
                        <div className="bg-white shadow-xl rounded-xl p-6 sm:p-8 h-full"> {/* h-full to match height if columns differ */}
                            <Form onSimulate={handleSimulate} isLoading={loading} />
                        </div>
                    </div>

                    {/* Column 2: Results (Chart, Table, Loading, Error) */}
                    {/* On large screens, this div takes 2/3 of the width. On smaller screens, it's full width. */}
                    <div className="lg:w-3/5 xl:w-2/3 space-y-6"> {/* space-y-6 for stacking items within this column */}
                        {loading && (
                            <div className="bg-white shadow-xl rounded-xl p-6 sm:p-8 text-center">
                                <p className="text-lg text-blue-600">Calculating your financial future... ⏳</p>
                            </div>
                        )}
                        {error && (
                            <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-md shadow-md" role="alert">
                                <p className="font-bold">Oops! Something went wrong.</p>
                                <p>{error}</p>
                            </div>
                        )}
                        {result && !loading && !error && (
                            <div className="bg-white shadow-xl rounded-xl p-6 sm:p-8 space-y-6">
                                <h2 className="text-2xl font-semibold text-gray-700 mb-4">Simulation Results</h2>
                                <Chart simData={result} />
                                <div className="overflow-x-auto mt-6">
                                    <table className="min-w-full table-auto border border-gray-300">
                                        <thead className="bg-gray-200">
                                        <tr>
                                            <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border-b-2 border-gray-300">Statistic</th>
                                            <th className="px-4 py-3 text-right text-sm font-semibold text-gray-700 border-b-2 border-gray-300">Value</th>
                                        </tr>
                                        </thead>
                                        <tbody className="divide-y divide-gray-200">
                                        {(Object.keys(result.final_stats) as Array<keyof SummaryStats>).map((key) => (
                                            <tr key={key} className="hover:bg-gray-50">
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-600 capitalize">{key.replace(/_/g, ' ')}</td>
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-800 text-right font-medium">
                                                    {formatCurrency(result.final_stats[key])}
                                                </td>
                                            </tr>
                                        ))}
                                        <tr className="hover:bg-gray-50">
                                            <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-600">Success Rate</td>
                                            <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-800 text-right font-medium">
                                                {(result.success_rate * 100).toFixed(1)}%
                                            </td>
                                        </tr>
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        )}
                        {/* Placeholder if no simulation run yet and not loading/error */}
                        {!result && !loading && !error && (
                            <div className="bg-white shadow-xl rounded-xl p-6 sm:p-8 text-center text-gray-500 h-full flex flex-col justify-center items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-12 w-12 mb-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1}>
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                                </svg>
                                <p className="text-lg">Run a simulation to see the results.</p>
                                <p className="text-sm">Your projected portfolio growth will appear here.</p>
                            </div>
                        )}
                    </div>
                </div>
                <footer className="text-center py-8 text-sm text-gray-500">
                    <p>Portfolio Simulator &copy; {new Date().getFullYear()}</p>
                </footer>
            </div>
        </div>
    );
}