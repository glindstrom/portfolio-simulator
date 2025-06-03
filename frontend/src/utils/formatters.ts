import React from 'react'; // Needed for React.ReactNode

// formatUSD formats a number as USD currency using en-US locale.
export const formatUSD = (value: number): string => {
    if (isNaN(value)) {
        return "N/A";
    }
    return value.toLocaleString('en-US', {
        style: "currency",
        currency: "USD",
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    });
};

// yAxisTickFormatter formats Y-axis ticks for large numbers.
export const yAxisTickFormatter = (value: number): string => {
    if (isNaN(value)) return "$0";
    if (Math.abs(value) >= 1e6) return `$${(value / 1e6).toFixed(1)}M`;
    if (Math.abs(value) >= 1e3) return `$${(value / 1e3).toFixed(0)}K`;
    return `$${value.toFixed(0)}`;
};

// Tooltip formatter function for Recharts.
// value: The value of the hovered data point for a series.
// name: The name of the hovered series.
export const tooltipFormatterFunc = (
    value: unknown,
    name: string,
): React.ReactNode | [React.ReactNode, React.ReactNode] => {

    if (name === '5th - 95th Percentile') {
        if (Array.isArray(value) && value.length === 2 &&
            typeof value[0] === 'number' && typeof value[1] === 'number') {
            const p5Value: number = value[0];
            const p95Value: number = value[1];
            return [`${formatUSD(p5Value)} - ${formatUSD(p95Value)}`, name];
        }
    }

    if (typeof value === 'number') {
        return [formatUSD(value), name];
    }

    if (value === undefined || value === null) {
        return ["N/A", name];
    }
    return [String(value), name];
};

// xAxisTickFormatter formats the X-axis ticks for years.
export const xAxisTickFormatter = (monthIndex: number): string => {
    const year = Math.floor(monthIndex / 12);
    return `${year}Y`; // e.g., "0Y", "1Y", "5Y"
};