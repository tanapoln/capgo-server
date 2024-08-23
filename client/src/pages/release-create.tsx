import { Button, Form, FormProps, Input, message, Radio, Select } from "antd";
import { createRelease, listBundles } from "../client/api";
import { useEffect, useState } from "react";
import { BundleResponse, Platform } from "../client/types";
import { useNavigate } from "react-router-dom";

type FieldType = {
	platform: string;
	app_id: string;
	version_name: string;
	version_code: string;
	builtin_bundle_id: string;
};

export default function ReleaseCreatePage() {
	const navigate = useNavigate();
	const [messageApi, contextHolder] = message.useMessage();
	const [bundles, setBundles] = useState<BundleResponse[]>([]);
	useEffect(() => {
		(async () => {
			const b = await listBundles();
			setBundles(b.data);
		})();
	}, []);

	const onFinish: FormProps<FieldType>["onFinish"] = async (values) => {
		try {
			await createRelease({
				platform: values.platform as Platform,
				app_id: values.app_id,
				version_name: values.version_name,
				version_code: values.version_code,
				builtin_bundle_id: values.builtin_bundle_id,
			});
			// navigate("/app");
		} catch (e) {
			messageApi.error(`${e}`);
		}
	};

	const onFinishFailed: FormProps<FieldType>["onFinishFailed"] = async (errorInfo) => {};

	return (
		<div>
			{contextHolder}
			<h1>Create Releases</h1>
			<div>
				<Form
					name="basic"
					labelCol={{ span: 8 }}
					wrapperCol={{ span: 16 }}
					style={{ maxWidth: 600 }}
					initialValues={{ remember: true }}
					onFinish={onFinish}
					onFinishFailed={onFinishFailed}
					autoComplete="off"
				>
					<Form.Item<FieldType>
						label="Platform"
						name="platform"
						rules={[{ required: true, message: "Please input value!" }]}
					>
						<Radio.Group>
							<Radio value="ios">iOS</Radio>
							<Radio value="android">Android</Radio>
						</Radio.Group>
					</Form.Item>

					<Form.Item<FieldType>
						label="App ID"
						name="app_id"
						rules={[{ required: true, message: "Please input value!" }]}
					>
						<Input />
					</Form.Item>

					<Form.Item<FieldType>
						label="Version Name"
						name="version_name"
						rules={[{ required: true, message: "Please input value!" }]}
					>
						<Input />
					</Form.Item>

					<Form.Item<FieldType>
						label="Version Code"
						name="version_code"
						rules={[{ required: true, message: "Please input value!" }]}
					>
						<Input />
					</Form.Item>

					<Form.Item<FieldType>
						label="Built-in Bundle ID"
						name="builtin_bundle_id"
						rules={[{ required: true, message: "Please input value!" }]}
					>
						<Select
							options={bundles.map((b) => ({
								value: b.id,
								label: `${b.version_name} - ${b.description}`,
							}))}
						/>
					</Form.Item>

					<Form.Item wrapperCol={{ offset: 8, span: 16 }}>
						<Button type="primary" htmlType="submit">
							Submit
						</Button>
					</Form.Item>
				</Form>
			</div>
		</div>
	);
}
