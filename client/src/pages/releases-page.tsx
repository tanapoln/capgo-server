import { NavLink } from "react-router-dom";
import type { ReleaseResponse } from "../client/types";
import { Button, Flex, Input, Table, type TableColumnsType } from "antd";
import { useReleases } from "../client/hooks";

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
		},
		{
			title: "Builtin Bundle ID",
			dataIndex: "builtin_bundle_id",
		},
		{
			title: "Active Bundle ID",
			dataIndex: "active_bundle_id",
		},
		{
			title: "Actions",
			dataIndex: "id",
			render: (_, record) => <NavLink to={`release/${record.id}/update`}>Update</NavLink>,
		},
	];

	return <Table dataSource={data} columns={columns} rowKey="id" />;
}
