export function UserNav({subdomain}: {subdomain: string}) {
    const userLink = `https://${subdomain}.quantinium.dev`;
    return (
        <nav className="flex justify-between flex-col mb-6">
            <div className="flex items-baseline gap-3 py-3">
                <h1 className="font-bold text-2xl">{subdomain}</h1>
                <span className="text-sm hover:underline">[<a href={userLink} className="px-1">{subdomain}.quantinium.dev</a>]</span>
            </div>
            <div className="flex gap-2 items-center text-[12px] lg:text-[15px]">
                <a className="hover:underline hover:underline-offset-2" href="/dashboard">home</a>
                <a className="hover:underline hover:underline-offset-2" href="/navigation">navigatoin</a>
                <a className="hover:underline hover:underline-offset-2" href="/posts">posts</a>
                <a className="hover:underline hover:underline-offset-2" href="/static">static pages</a>
                <a className="hover:underline hover:underline-offset-2" href="/theme">themes</a>
                <a className="hover:underline hover:underline-offset-2" href="/setting">settings</a>
            </div>
        </nav>
    )
}
