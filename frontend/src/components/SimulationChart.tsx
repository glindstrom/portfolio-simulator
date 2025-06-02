import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';

interface ChartProps {
    data: { month: number; median: number; p5: number; p95: number }[];
}

export default function SimulationChart({ data }: ChartProps) {
    return (
        <div className="w-full h-96 bg-white rounded-xl shadow p-4">
            <h2 className="text-lg font-semibold mb-2">Portfolio Value Over Time</h2>
            <ResponsiveContainer width="100%" height="90%">
                <LineChart data={data}>
                    <XAxis dataKey="month" />
                    <YAxis />
                    <Tooltip />
                    <Line type="monotone" dataKey="median" stroke="#3b82f6" strokeWidth={2} dot={false} />
                    <Line type="monotone" dataKey="p5" stroke="#f87171" strokeDasharray="3 3" dot={false} />
                    <Line type="monotone" dataKey="p95" stroke="#34d399" strokeDasharray="3 3" dot={false} />
                </LineChart>
            </ResponsiveContainer>
        </div>
    );
}
