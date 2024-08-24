import { Button, Card, Col, DatePicker, Flex, Form, message, Popconfirm, Row, Select } from "antd";
import type { Dayjs } from "dayjs";
import dayjs from "dayjs";
import { useNavigate, useParams } from "react-router-dom";
import {
	useBundles,
	useDeleteReleaseMutation,
	useRelease,
	useSetReleaseActiveBundleMutation,
	useUpdateReleaseMutation,
} from "../client/hooks";
import { ReleaseResponse } from "../client/types";
import { BundlePopover } from "../ui/bundle";

export default function ReleaseUpdatePage() {
	const { releaseId } = useParams();
	const { data: release, isLoading, error } = useRelease(releaseId!);

	return (
		<div style={{ padding: "24px 0" }}>
			<h1>Update Releases</h1>
			{isLoading && <div>Loading...</div>}
			{error && <div>Error: {error.message}</div>}
			{release && (
				<Flex vertical gap="large" style={{ maxWidth: "650px" }}>
					<ReleaseInfo release={release} />
					<SetReleaseDateForm release={release} />
					<SetActiveBundleForm release={release} />
					<DeleteReleaseForm release={release} />
				</Flex>
			)}
		</div>
	);
}

function ReleaseInfo({ release }: { release: ReleaseResponse }) {
	return (
		<Card title="Release Information" style={{ minWidth: "300px" }}>
			<Flex vertical gap="small">
				<Row>
					<Col span={6} style={{ fontWeight: 500 }}>
						Platform
					</Col>
					<Col span={18}>{release.platform}</Col>
				</Row>
				<Row>
					<Col span={6} style={{ fontWeight: 500 }}>
						App ID
					</Col>
					<Col span={18}>{release.app_id}</Col>
				</Row>
				<Row>
					<Col span={6} style={{ fontWeight: 500 }}>
						Version Name
					</Col>
					<Col span={18}>{release.version_name}</Col>
				</Row>
				<Row>
					<Col span={6} style={{ fontWeight: 500 }}>
						Version Code
					</Col>
					<Col span={18}>{release.version_code}</Col>
				</Row>
				<Row>
					<Col span={6} style={{ fontWeight: 500 }}>
						Release Date
					</Col>
					<Col span={18}>{dayjs(release.release_date).format("YYYY-MM-DD HH:mm:ss")}</Col>
				</Row>
				<Row>
					<Col span={6} style={{ fontWeight: 500 }}>
						Builtin Bundle ID
					</Col>
					<Col span={18}>
						<BundlePopover bundle_id={release.builtin_bundle_id} />
					</Col>
				</Row>
				<Row>
					<Col span={6} style={{ fontWeight: 500 }}>
						Active Bundle ID
					</Col>
					<Col span={18}>
						{release.active_bundle_id && <BundlePopover bundle_id={release.active_bundle_id} />}
					</Col>
				</Row>
				<Row>
					<Col span={6} style={{ fontWeight: 500 }}>
						Created At
					</Col>
					<Col span={18}>{release.created_at}</Col>
				</Row>
			</Flex>
		</Card>
	);
}

function SetReleaseDateForm({ release }: { release: ReleaseResponse }) {
	const { trigger } = useUpdateReleaseMutation();

	const onSubmit = (values: { release_date: Dayjs }) => {
		const data = { release_id: release.id, release_date: values.release_date.toISOString() };
		trigger(data);
		message.success("Release date updated");
	};

	return (
		<Card title="Set Release Date">
			<Form onFinish={onSubmit}>
				<Form.Item
					label="Release Date"
					name="release_date"
					rules={[{ required: true, message: "Date is required" }]}
				>
					<DatePicker />
				</Form.Item>
				<Button type="primary" htmlType="submit">
					Save
				</Button>
			</Form>
		</Card>
	);
}

function SetActiveBundleForm({ release }: { release: ReleaseResponse }) {
	const { trigger } = useSetReleaseActiveBundleMutation();
	const { data: bundles, isLoading, error } = useBundles();

	const onSubmit = (values: { bundle_id: string }) => {
		const data = { release_id: release.id, bundle_id: values.bundle_id };
		trigger(data);
		message.success("Active bundle is updated");
	};

	return (
		<Card title="Set active bundle">
			<Form onFinish={onSubmit}>
				<Form.Item
					label="Bundle ID"
					name="bundle_id"
					rules={[{ required: true, message: "Please input value!" }]}
				>
					{isLoading && <Select options={[{ value: "", label: "Loading..." }]} />}
					{error && <div>Error: {error.message}</div>}
					{bundles && (
						<Select
							options={bundles.data.map((b) => ({
								value: b.id,
								label: `${b.app_id} - ${b.version_name}${
									b.description !== "" ? ` [${b.description}]` : ""
								}`,
							}))}
						/>
					)}
				</Form.Item>
				<Button type="primary" htmlType="submit">
					Save
				</Button>
			</Form>
		</Card>
	);
}

function DeleteReleaseForm({ release }: { release: ReleaseResponse }) {
	const navigate = useNavigate();
	const { trigger } = useDeleteReleaseMutation();

	const onSubmit = () => {
		trigger({ release_id: release.id });
		message.success("Release deleted");
		navigate("/app")
	};

	return (
		<Card title={<span style={{ color: "red" }}>Delete Release</span>}>
			<Popconfirm
				title="Confirm delete?"
				okText="Delete"
				onConfirm={onSubmit}
			>
				<Button type="dashed" danger>
					Delete
				</Button>
			</Popconfirm>
		</Card>
	);
}
