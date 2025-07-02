"use client"

import { useState } from "react";
import dynamic from "next/dynamic";
import "react-quill-new/dist/quill.snow.css";
import { z } from "zod";

const modules = {
    toolbar: [
        [{ header: [1, 2, 3, 4, 5, false] }],
        ["bold", "italic", "underline"],
        [{ color: [] }, { background: [] }],
        ["blockquote", "code-block"],
        [{ list: "ordered" }, { list: "bullet" }],
        [{ align: [] }],
        [{ size: ["small", false, "large", "huge"] }],
        ["link", "image"],
    ],
};

const formats = [
    "header",
    "bold",
    "color",
    "background",
    "italic",
    "underline",
    "blockquote",
    "code-block",
    "list",
    "align",
    "size",
    "link",
    "image",
];

const ReactQuill = dynamic(() => import("react-quill-new"), {
    ssr: false
})

export default function DashboardPage() {
    const [value, setValue] = useState("");
    const [metadataVisibility, setMetadataVisibility] = useState(false);
    const [title, setTitle] = useState("");
    const [faviconUrl, setFavionUrl] = useState("");
    const [metadataDescription, setMetadataDescription] = useState("");
    const [metadataImage, setMetadataImage] = useState("");

    function toggleMetadata() {
        setMetadataVisibility(!metadataVisibility);
    }

    function onFaviconChange(e) {
        e.preventDefault()
        setFavionUrl(e.target.value);
    }

    function onMetadataDescription(e) {
        e.prevemtDefault();
        setMetadataDescription(e.target.value)
    }

    function onMetadataImage(e) {
        e.preventDefault();
        setMetadataImage(e.target.value)
    }

    function handlePublish() {
        console.log(value)
        console.log(title)
        console.log(faviconUrl)
        console.log(metadataDescription);
        console.log(metadataImage);
    }

    function onTitleChange(e) {
        e.preventDefault()
        setTitle(e.target.value);
    }

    return (
        <>
            <div className="font-semibold text-xl">Home</div>
            <div className="flex gap-3">
                <button onClick={handlePublish} className="border px-1 border-black hover:bg-gray-100 rounded my-3 text-[15px]">publish</button>
                <button className="border px-1 border-black hover:bg-gray-100 rounded my-3 text-[15px]">preview</button>
                <button onClick={toggleMetadata} className="border px-1 border-black hover:bg-gray-100 rounded my-3 text-[15px]">metadata</button>
            </div>

            <div
                className={`transition-all duration-500 ease-in-out mb-3 ${metadataVisibility ? 'max-h-0 opacity-0 overflow-hidden' : 'max-h-[1000px] opacity-100'
                    }`}
            >
                <div className="flex flex-col gap-3 text-[12px]">
                    <label className="flex gap-3"><span>Title:</span>
                        <input type="text" name="title" className="border-2 border-black title" value={title} onChange={onTitleChange} />
                    </label>
                    <label className="flex gap-3"><span>Favicon Url:</span>
                        <input type="url" name="favicon" className="p-0 border-2 border-black favicon" value={faviconUrl} onChange={onFaviconChange} />
                    </label>
                    <label className="flex gap-3"><span>Metadata Description:</span>
                        <input type="text" name="metadata_description" className="p-0 border-2 border-black metadata_description" value={metadataDescription} onChange={onMetadataDescription} />
                    </label>
                    <label className="flex gap-3"><span>Metadata Url:</span>
                        <input type="text" name="metadata_url" className="p-0 border-2 border-black metadata_url" value={metadataImage} onChange={onMetadataImage} />
                    </label>
                </div>

            </div >

            <ReactQuill
                theme="snow"
                value={value}
                onChange={setValue}
                modules={modules}
                formats={formats}
                placeholder="Write something awesome..."
            />
        </>
    );
}
