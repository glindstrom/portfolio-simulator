import { useState } from "react";
import Form from "./components/Form";
import ResultsDisplay from "./components/ResultsDisplay"; // Import the new component
import type { SimulationResponse, Params } from "./types"; // SummaryStats might not be needed directly in App.tsx now

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
                // Handle HTTP errors (like 4xx, 5xx) directly
                let errorMessage = `HTTP error! Status: ${res.status}`;
                try {
                    const errorData = await res.json(); // Try to parse error response from backend
                    errorMessage = errorData.detail || errorData.message || errorMessage;
                } catch (jsonParseError) {
                    // If error response is not JSON, or parsing fails, stick with the HTTP status error
                    console.warn("Could not parse error response JSON:", jsonParseError);
                }
                console.error("API Error:", errorMessage);
                setError(errorMessage);
                return;
            }

            // If res.ok is true, proceed to parse the successful response
            const data: SimulationResponse = await res.json();
            console.log("Received simulation result:", data);
            setResult(data);

        } catch (err: any) {
            console.error("Simulation request processing failed:", err);
            setError(err.message || "An unexpected error occurred.");
        } finally {
            setLoading(false);
        }
    };

    // formatCurrency is now primarily used within ResultsDisplay.tsx.
    // If not used elsewhere in App.tsx, it can be removed from here.
    // For now, let's assume it might be used by an older Chart.tsx if that's still around,
    // or it can be safely removed if ResultsDisplay.tsx is the sole consumer.

    return (
        <div className="min-h-screen bg-gray-100 p-4 sm:p-6 font-sans">
            <div className="max-w-7xl mx-auto space-y-8">
                <header className="text-center py-6">
                    <h1 className="text-3xl sm:text-4xl font-bold text-gray-800">
                        📈 Portfolio Growth Simulator
                    </h1>
                </header>

                <div className="lg:flex lg:gap-8 space-y-8 lg:space-y-0">
                    <div className="lg:w-2/5 xl:w-1/3">
                        <div className="bg-white shadow-xl rounded-xl p-6 sm:p-8 h-full">
                            <Form onSimulate={handleSimulate} isLoading={loading} />
                        </div>
                    </div>

                    <div className="lg:w-3/5 xl:w-2/3 space-y-6">
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
                            // Use the new ResultsDisplay component here
                            <ResultsDisplay resultData={result} />
                        )}
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