"use client"

import Link from "next/link";

type Post = {
    title: string
    description: string
    slug: string
    imageUrl: string
    date: string
    tags: string[]
    isPublished: boolean
}

const posts: Post[] = [
    {
        title: "Getting Started with TypeScript",
        description: "A beginner's guide to understanding and using TypeScript in your projects.",
        slug: "getting-started-with-typescript",
        imageUrl: "https://source.unsplash.com/featured/?typescript,code",
        date: "2025-06-15",
        tags: ["typescript", "javascript"],
        isPublished: true
    },
    {
        title: "Mastering Tailwind CSS for Responsive Design",
        description: "Explore how Tailwind CSS can help you build responsive and modern UI components.",
        slug: "mastering-tailwind-css",
        imageUrl: "https://source.unsplash.com/featured/?tailwindcss,ui",
        date: "2025-06-22",
        tags: ["css", "tailwind"],
        isPublished: true
    },
    {
        title: "Understanding React Server Components",
        description: "A deep dive into React Server Components and their role in modern web development.",
        slug: "understanding-react-server-components",
        imageUrl: "https://source.unsplash.com/featured/?react,server",
        date: "2025-06-25",
        tags: ["react", "server-components"],
        isPublished: false
    },
    {
        title: "Building REST APIs with Express.js",
        description: "Learn how to create scalable and maintainable REST APIs using Express.",
        slug: "building-rest-apis-with-express",
        imageUrl: "https://source.unsplash.com/featured/?express,api",
        date: "2025-05-30",
        tags: ["nodejs", "express"],
        isPublished: true
    },
    {
        title: "Docker Essentials for Developers",
        description: "An introduction to Docker, containers, and how to use them in your development workflow.",
        slug: "docker-essentials-for-developers",
        imageUrl: "https://source.unsplash.com/featured/?docker,containers",
        date: "2025-06-10",
        tags: ["docker", "devops"],
        isPublished: true
    }
];

export default function Posts() {
    return (
        <>
            <div>
                <Link href="/posts/new" className="border-black border px-1 rounded hover:bg-zinc-100">new post</Link>
            </div>

            <div className="my-3">
                {posts.map((post) => (
                    <div key={post.title} className="border-b py-3">
                        <h1 className="font-semibold text-xl">{post.title}</h1>
                        <div className="flex gap-3 text-sm pb-3">
                            <span className="border-r border-black pr-3">{post.tags.join(", ")}</span>
                            <span>{post.date}</span>
                        </div>

                        <div className="pb-2">
                            {post.description}
                        </div>
                        <Link href="https://github.com" className="hover:underline hover:underline-offset-2">more...</Link>
                    </div>
                ))}
            </div>
        </>
    );
}
