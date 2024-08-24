import { useParams } from "react-router-dom";
import { useRelease } from "../client/hooks";

export default function ReleaseUpdatePage() {
	const { releaseId } = useParams();
	const { data: release, isLoading, error } = useRelease(releaseId!);

	return (
		<div>
			<h1>Update Releases</h1>
			{isLoading && <div>Loading...</div>}
			{error && <div>Error: {error.message}</div>}
			{release && <div>Release: {release.id}</div>}
		</div>
	);
}
