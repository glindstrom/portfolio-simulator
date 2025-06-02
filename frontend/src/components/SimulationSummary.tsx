interface SummaryProps {
    finalValues: {
        percentile: string;
        value: number;
    }[];
}

export default function SimulationSummary({ finalValues }: SummaryProps) {
    return (
        <div className="mt-6 bg-white rounded-xl shadow p-4">
            <h2 className="text-lg font-semibold mb-2">Summary of Final Portfolio Values</h2>
            <table className="w-full text-sm text-left border-t border-gray-200">
                <thead>
                <tr>
                    <th className="py-2">Percentile</th>
                    <th className="py-2">Final Value</th>
                </tr>
                </thead>
                <tbody>
                {finalValues.map((row) => (
                    <tr key={row.percentile} className="border-t border-gray-100">
                        <td className="py-2">{row.percentile}</td>
                        <td className="py-2">{row.value.toLocaleString()}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
}
