// Defines the shape of data points for the chart.
export type ChartPoint = {
    month: number;    // X-axis: 0-indexed month number for internal data mapping
    median: number;   // Y-axis: median portfolio value
    p5_p95_range: [number, number]; // Range for 5th-95th percentile [p5, p95]
    // Individual p5 and p95 can be added if needed for other chart elements,
    // but p5_p95_range is used by the main Area component.
};

// percentile calculates the p-th percentile value in a sorted numeric array.
export function percentile(sortedArr: number[], p: number): number {
    if (!sortedArr || sortedArr.length === 0) return 0;
    const N = sortedArr.length;
    if (N === 1) return sortedArr[0]; // Single element is all percentiles

    const index = (p / 100) * (N - 1);
    const lower = Math.floor(index);
    const upper = Math.ceil(index);

    if (lower === upper) return sortedArr[lower]; // Exact match or p=0 or p=100

    return sortedArr[lower] + (sortedArr[upper] - sortedArr[lower]) * (index - lower);
}

// calculateChartData transforms raw simulation paths into data suitable for the chart.
export function calculateChartData(paths: number[][]): ChartPoint[] {
    if (!paths || paths.length === 0 || !paths[0] || paths[0].length === 0) {
        return [];
    }
    const numMonthsInPath = paths[0].length;
    const chartData: ChartPoint[] = [];

    for (let i = 0; i < numMonthsInPath; i++) { // Loop through each month/period index
        const valuesAtThisPeriod: number[] = [];
        for (let j = 0; j < paths.length; j++) { // Loop through each simulation path
            if (paths[j] && typeof paths[j][i] === 'number') {
                valuesAtThisPeriod.push(paths[j][i]);
            }
        }
        valuesAtThisPeriod.sort((a, b) => a - b);

        if (valuesAtThisPeriod.length > 0) {
            const p5Val = percentile(valuesAtThisPeriod, 5);
            const p95Val = percentile(valuesAtThisPeriod, 95);
            chartData.push({
                month: i, // 0-indexed month for dataKey
                median: percentile(valuesAtThisPeriod, 50),
                p5_p95_range: [p5Val, p95Val],
            });
        }
    }
    return chartData;
}

// generateXAxisYearTicks creates an array of ticks for the X-axis based on the total months in the data.
export function generateXAxisYearTicks(totalMonthsInData: number): number[] {
    const xTicks: number[] = [];
    if (totalMonthsInData <= 0) {
        return xTicks;
    }

    const tickSet = new Set<number>();
    tickSet.add(0); // Month 0 (represents start of Year 0 or beginning)

    const maxYearDisplay = Math.floor((totalMonthsInData - 1) / 12);

    // Add Year 1 (Month 12) if distinct and simulation is long enough
    if (maxYearDisplay >= 1 && totalMonthsInData > 12) {
        tickSet.add(12);
    }

    // Add ticks at 5-year intervals
    for (let year = 5; year <= maxYearDisplay; year += 5) {
        // Avoid duplicating Year 1 if it's also a 5-year mark (unlikely with year=1 already added)
        // or if maxYearDisplay is small.
        if (year * 12 < totalMonthsInData) { // Ensure tick is within data range
            tickSet.add(year * 12);
        }
    }

    // Ensure the last data point (month index) is always a tick
    if (totalMonthsInData > 1) {
        tickSet.add(totalMonthsInData - 1);
    }

    return Array.from(tickSet).sort((a, b) => a - b);
}