// src/types/index.ts

// Represents one portfolio asset and its weight
export type PortfolioItem = {
    ticker: string;
    weight: number;
};

// Parameters sent to the backend simulation API
export type Params = {
    portfolio: PortfolioItem[];
    initialValue: number;
    periods: number;
    simulations: number;
    method: "normal" | "bootstrap";
    withdrawal: number;
    inflation: number;
};

// Statistics returned after simulation
export type SummaryStats = {
    mean: number;
    median: number;
    min: number;
    max: number;
};

// Full response from the backend simulation API
export type SimulationResponse = {
    paths: number[][];
    finalStats: SummaryStats;
    successRate: number; // between 0 and 1
    simulatedCAGR: number; // Compound Annual Growth Rate of the simulated paths, including withdrawals
};
