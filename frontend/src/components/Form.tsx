// src/components/Form.tsx
import React, { useState } from "react"; // Import React itself
import type { Params, PortfolioItem } from "../types";

// Updated Props type
type Props = {
    onSimulate: (params: Params) => Promise<void>; // Matched Promise<void> from App.tsx
    isLoading: boolean;                             // Added isLoading
};

export default function Form({ onSimulate, isLoading }: Props) { // Destructure isLoading
    const [portfolio, setPortfolio] = useState<PortfolioItem[]>([
        { ticker: "AAPL", weight: 0.6 },
        { ticker: "MSFT", weight: 0.4 },
    ]);
    const [initialValue, setInitialValue] = useState(10000);
    const [years, setYears] = useState(40);
    const [simulations, setSimulations] = useState(500);
    const [method, setMethod] = useState<"normal" | "bootstrap">("bootstrap");
    // Assuming withdrawal is a rate (e.g., 0.04 for 4%) not an absolute amount per period
    const [withdrawalRate, setWithdrawalRate] = useState(0.04);
    const [inflation, setInflation] = useState(0.02);

    const handleSubmit = async (e: React.FormEvent) => { // Made handleSubmit async
        e.preventDefault();
        await onSimulate({ // Awaited the onSimulate call
            portfolio,
            initialValue,
            periods: years * 12, // Convert years to months
            simulations,
            method,
            withdrawal: withdrawalRate,
            inflation,
        });
    };

    const handlePortfolioChange = (index: number, field: keyof PortfolioItem, value: string) => {
        const newPortfolio = [...portfolio];
        if (field === "weight") {
            // Ensure weight is a number between 0 and 1 if it's a percentage
            const weightValue = parseFloat(value);
            newPortfolio[index].weight = isNaN(weightValue) ? 0 : weightValue; // Basic NaN check
        } else if (field === "ticker") {
            newPortfolio[index].ticker = value.toUpperCase(); // Standardize ticker to uppercase
        }
        setPortfolio(newPortfolio);
    };

    const addAsset = () => {
        // Add validation to prevent adding too many empty assets or ensure weights sum to 1 eventually
        setPortfolio([...portfolio, { ticker: "", weight: 0 }]);
    };

    const removeAsset = (index: number) => {
        setPortfolio(portfolio.filter((_, i) => i !== index));
    };

    // Basic validation for weights (sum to 1) - for display/UX purposes
    const totalWeight = portfolio.reduce((sum, asset) => sum + (asset.weight || 0), 0);
    const isWeightInvalid = parseFloat(totalWeight.toFixed(4)) !== 1.0;


    return (
        <form onSubmit={handleSubmit} className="space-y-6 p-4 border border-gray-200 rounded-lg shadow-sm bg-white h-full flex flex-col">
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold text-gray-700">Simulation Parameters</h2>
            </div>

            <div>
                <h3 className="text-lg font-medium text-gray-600 mb-2">Portfolio Assets:</h3>
                {portfolio.map((asset, index) => (
                    <div key={index} className="flex gap-2 items-center mb-2 p-2 border rounded-md bg-gray-50">
                        <input
                            type="text"
                            placeholder="Ticker (e.g., SPY)"
                            value={asset.ticker}
                            onChange={(e) => handlePortfolioChange(index, "ticker", e.target.value)}
                            className="border p-2 rounded w-2/5 focus:ring-indigo-500 focus:border-indigo-500"
                            disabled={isLoading}
                        />
                        <input
                            type="number"
                            placeholder="Weight (e.g., 0.6)"
                            step="0.01"
                            min="0"
                            max="1"
                            value={asset.weight}
                            onChange={(e) => handlePortfolioChange(index, "weight", e.target.value)}
                            className="border p-2 rounded w-2/5 focus:ring-indigo-500 focus:border-indigo-500"
                            disabled={isLoading}
                        />
                        <button
                            type="button"
                            onClick={() => removeAsset(index)}
                            className="text-red-500 hover:text-red-700 px-2 py-1 rounded disabled:opacity-50"
                            disabled={isLoading || portfolio.length <= 1} // Prevent removing the last asset
                        >
                            ✕
                        </button>
                    </div>
                ))}
                <div className="flex items-center gap-4 mt-2">
                    <button
                        type="button"
                        onClick={addAsset}
                        className="bg-blue-500 text-white px-3 py-1.5 rounded hover:bg-blue-600 text-sm disabled:bg-gray-400"
                        disabled={isLoading}
                    >
                        + Add Asset
                    </button>
                    {isWeightInvalid && (
                        <p className="text-sm text-red-600">Total weight must be 1.0 (current: {totalWeight.toFixed(4)})</p>
                    )}
                </div>
            </div>


            <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4 pt-4 border-t mt-6">
                <div>
                    <label htmlFor="initialValue" className="block text-sm font-medium text-gray-700">Initial Value ($)</label>
                    <input
                        id="initialValue"
                        type="number"
                        min="0"
                        value={initialValue}
                        onChange={(e) => setInitialValue(parseFloat(e.target.value))}
                        className="mt-1 border p-2 rounded w-full focus:ring-indigo-500 focus:border-indigo-500"
                        disabled={isLoading}
                    />
                </div>

                <div>
                    <label htmlFor="years" className="block text-sm font-medium text-gray-700">Investment Horizon (Years)</label>
                    <input
                        id="years"
                        type="number"
                        min="1"
                        value={years}
                        onChange={(e) => setYears(parseInt(e.target.value))}
                        className="mt-1 border p-2 rounded w-full focus:ring-indigo-500 focus:border-indigo-500"
                        disabled={isLoading}
                    />
                </div>

                <div>
                    <label htmlFor="simulations" className="block text-sm font-medium text-gray-700">Number of Simulations</label>
                    <input
                        id="simulations"
                        type="number"
                        min="1"
                        value={simulations}
                        onChange={(e) => setSimulations(parseInt(e.target.value))}
                        className="mt-1 border p-2 rounded w-full focus:ring-indigo-500 focus:border-indigo-500"
                        disabled={isLoading}
                    />
                </div>

                <div>
                    <label htmlFor="method" className="block text-sm font-medium text-gray-700">Simulation Method</label>
                    <select
                        id="method"
                        value={method}
                        onChange={(e) => setMethod(e.target.value as "normal" | "bootstrap")}
                        className="mt-1 border p-2 rounded w-full focus:ring-indigo-500 focus:border-indigo-500 bg-white"
                        disabled={isLoading}
                    >
                        <option value="normal">Normal Distribution</option>
                        <option value="bootstrap">Historical Bootstrap</option>
                    </select>
                </div>

                <div>
                    <label htmlFor="withdrawalRate" className="block text-sm font-medium text-gray-700">Annual Withdrawal Rate (%)</label>
                    <input
                        id="withdrawalRate"
                        type="number"
                        step="0.001" // e.g. for 0.04 (4%) or 0.035 (3.5%)
                        min="0"
                        max="1" // Assuming rate is between 0 and 1
                        value={withdrawalRate}
                        onChange={(e) => setWithdrawalRate(parseFloat(e.target.value))}
                        className="mt-1 border p-2 rounded w-full focus:ring-indigo-500 focus:border-indigo-500"
                        placeholder="e.g., 0.04 for 4%"
                        disabled={isLoading}
                    />
                </div>

                <div>
                    <label htmlFor="inflation" className="block text-sm font-medium text-gray-700">Annual Inflation Rate (%)</label>
                    <input
                        id="inflation"
                        type="number"
                        step="0.001"
                        min="0"
                        value={inflation}
                        onChange={(e) => setInflation(parseFloat(e.target.value))}
                        className="mt-1 border p-2 rounded w-full focus:ring-indigo-500 focus:border-indigo-500"
                        placeholder="e.g., 0.02 for 2%"
                        disabled={isLoading}
                    />
                </div>
            </div>

            <div className="pt-6 border-t mt-6 flex-grow flex flex-col justify-end">
                <button
                    type="submit"
                    disabled={isLoading || isWeightInvalid} // Also disable if weights are invalid
                    className="w-full bg-green-600 text-white px-6 py-3 rounded-md text-lg font-semibold hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:bg-gray-400 disabled:cursor-not-allowed"
                >
                    {isLoading ? 'Simulating...' : 'Run Simulation'}
                </button>
                {isWeightInvalid && !isLoading && (
                    <p className="text-sm text-red-600 text-center mt-2">Please ensure total portfolio weight sums to 1.0 before simulating.</p>
                )}
            </div>
        </form>
    );
}