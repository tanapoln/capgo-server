import { Popover, Space } from "antd";
import { useBundle } from "../client/hooks";
import { BundleResponse } from "../client/types";

export function BundlePopover({ bundle_id }: { bundle_id: string }) {
	const { data: bundle, isLoading, error } = useBundle(bundle_id);

	const title = bundle ? `${bundle.app_id} - ${bundle.version_name}` : "";

	return (
		<>
			{isLoading && <div>Loading...</div>}
			{error && <div>Error: {error.message}</div>}
			{bundle && (
				<Popover title={title} content={<BundlePopoverContent bundle={bundle} />}>
					<span style={{ cursor: "pointer" }}>{title}</span>
				</Popover>
			)}
		</>
	);
}

function BundlePopoverContent({ bundle }: { bundle: BundleResponse }) {
	return (
		<Space direction="vertical">
			<div>App ID: {bundle.app_id}</div>
			<div>Version Name: {bundle.version_name}</div>
			<div>Description: {bundle.description}</div>
			<div>Checksum: {bundle.crc_checksum}</div>
			<div>Created Date: {bundle.created_at}</div>
			<div>
				Download URL:{" "}
				<a href={bundle.public_download_url} target="_blank" rel="noreferrer">
					Download
				</a>
			</div>
		</Space>
	);
}
