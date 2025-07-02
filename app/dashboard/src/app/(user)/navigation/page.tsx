"use client"

import { ChangeEvent, useState } from "react";

type Field = {
    key: string;
    value: string;
};

export default function Navigation() {
    const [fields, setFields] = useState<Field[]>([]);
    const [data, setData] = useState<Record<string, string>>({
        home: "/"
    });

    const handleFields = () => {
        setFields([...fields, { key: "", value: "" }]);
    };

    const handleChange = (
        index: number,
        type: "key" | "value",
        event: ChangeEvent<HTMLInputElement>
    ) => {
        const updatedFields = [...fields];
        updatedFields[index][type] = event.target.value;
        setFields(updatedFields);
    };

    const handleSave = () => {
        const newData: Record<string, string> = { ...data };

        fields.forEach(({ key, value }) => {
            if (key.trim()) {
                newData[key.trim()] = value;
            }
        });

        setData(newData);
        setFields([]);
    };

    function handlePublish() {
        // send the json to the backend with a fetch request
    }

    return (
        <div className="space-y-4">
            <h1 className="font-semibold text-xl">Navigation</h1>

            <div className="flex gap-3">
                <button
                    onClick={handleSave}
                    className="border rounded border-black px-1 hover:bg-zinc-100"
                >
                    save
                </button>
                <button
                    onClick={handlePublish}
                    className="border rounded border-black px-1 hover:bg-zinc-100"
                >
                    publish
                </button>
            </div>

            {fields.map((field, index) => (
                <div key={index} className="flex space-x-2">
                    <input
                        type="text"
                        placeholder="Key"
                        value={field.key}
                        onChange={(e) => handleChange(index, "key", e)}
                        className="border px-2 py-1 rounded"
                    />
                    <input
                        type="text"
                        placeholder="Value"
                        value={field.value}
                        onChange={(e) => handleChange(index, "value", e)}
                        className="border px-2 py-1 rounded"
                    />
                </div>
            ))}

            <button
                onClick={handleFields}
                className="border border-black h-6 w-6 rounded hover:bg-zinc-100 flex items-center justify-center"
            >
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                    className="w-4 h-4"
                >
                    <path
                        d="M6 12H18M12 6V18"
                        stroke="#000000"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                    />
                </svg>
            </button>

            <span className="my-3">preview: </span>
            <pre className="bg-gray-100 p-2 rounded text-sm">
                {JSON.stringify(data, null, 2)}
            </pre>
        </div>
    );
}
