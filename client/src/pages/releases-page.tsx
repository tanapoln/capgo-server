import { Button, Flex, Input, Popconfirm, Table, type TableColumnsType } from "antd";
import dayjs from "dayjs";
import { useNavigate } from "react-router-dom";
import { useDeleteReleaseMutation, useReleases } from "../client/hooks";
import type { ReleaseResponse } from "../client/types";
import { BundlePopover } from "../ui/bundle";

export default function ReleasesPage() {
	const { data: releases, isLoading, error } = useReleases();

	return (
		<div style={{ padding: "24px 0" }}>
			<h1>Releases</h1>
			{isLoading && <div>Loading...</div>}
			{error && <div>Error: {error.message}</div>}
			{releases && <ReleaseTable data={releases.data} />}
		</div>
	);
}

function ReleaseTable({ data }: { data: ReleaseResponse[] }) {
	const navigate = useNavigate();
	const { trigger: deleteRelease } = useDeleteReleaseMutation();

	const columns: TableColumnsType<ReleaseResponse> = [
		{
			title: "ID",
			dataIndex: "id",
			key: "id",
		},
		{
			title: "Platform",
			dataIndex: "platform",
			filters: [
				{
					text: "iOS",
					value: "ios",
				},
				{
					text: "Android",
					value: "android",
				},
			],
			onFilter: (value, record) => record.platform === value,
		},
		{
			title: "App ID",
			dataIndex: "app_id",
			filters: data.map((r) => ({
				text: r.app_id,
				value: r.app_id,
			})),
		},
		{
			title: "Version Name",
			dataIndex: "version_name",
			filterDropdown({ selectedKeys, setSelectedKeys, confirm }) {
				return (
					<div style={{ padding: 8 }} onKeyDown={(e) => e.stopPropagation()}>
						<Input
							placeholder="Search version name"
							value={selectedKeys[0]}
							onChange={(e) => setSelectedKeys(e.target.value ? [e.target.value] : [])}
							onPressEnter={() => confirm()}
							style={{ marginBottom: 8, display: "block" }}
						/>
						<Flex align="center" justify="flex-end" gap={8}>
							<Button
								type="default"
								size="small"
								onClick={() => {
									setSelectedKeys([]);
									confirm();
								}}
							>
								Clear
							</Button>
							<Button type="primary" size="small" onClick={() => confirm()}>
								Search
							</Button>
						</Flex>
					</div>
				);
			},

			onFilter(value, record) {
				return record.version_name.includes(value as string);
			},
		},
		{
			title: "Version Code",
			dataIndex: "version_code",
			sorter: (a, b) => a.version_code.localeCompare(b.version_code),
		},
		{
			title: "Release Date",
			dataIndex: "release_date",
			sorter: (a, b) => (a.release_date ?? "").localeCompare(b.release_date ?? ""),
			render: (_, record) => record.release_date ? dayjs(record.release_date).format("YYYY-MM-DD HH:mm:ss") : "",
		},
		{
			title: "Builtin Bundle ID",
			dataIndex: "builtin_bundle_id",
			render: (_, record) => <BundlePopover bundle_id={record.builtin_bundle_id} />,
		},
		{
			title: "Active Bundle ID",
			dataIndex: "active_bundle_id",
			render: (_, record) =>
				record.active_bundle_id ? <BundlePopover bundle_id={record.active_bundle_id} /> : "[N/A]",
		},
		{
			title: "Actions",
			dataIndex: "id",
			align: "center",
			render: (_, record) => (
				<div>
					<Button type="link" onClick={() => navigate(`release/${record.id}/update`)}>
						Update
					</Button>
					<Popconfirm
						title="Delete this release?"
						okText="Delete"
						onConfirm={() => deleteRelease({ release_id: record.id })}
					>
						<Button type="dashed" danger>
							Delete
						</Button>
					</Popconfirm>
				</div>
			),
		},
	];

	return <Table dataSource={data} columns={columns} rowKey="id" />;
}
