import { Await } from "react-router-dom";
import { listReleases } from "../client/api";
import { Suspense } from "react";
import { ListAllReleasesResponse } from "../client/types";

export default function ReleasesPage() {
	const releases = listReleases();

	return (
		<div>
			<h1>Releases</h1>
			<Suspense fallback={<div>Loading...</div>}>
				<Await
					resolve={releases}
					children={(resolved: ListAllReleasesResponse) => (
						<div>
							{resolved.data.map((release) => (
								<div key={release.id}>{release.version_name}</div>
							))}
						</div>
					)}
				/>
			</Suspense>
		</div>
	);
}
