import { UploadOutlined } from "@ant-design/icons";
import { Button, Form, type FormProps, Input, message, Upload, UploadFile } from "antd";
import { useNavigate } from "react-router-dom";
import { useUploadBundleMutation } from "../client/hooks";

type FieldType = {
	bundle: UploadFile[];
	app_id: string;
	version_name: string;
	description?: string;
};

export default function BundleUploadPage() {
	const navigate = useNavigate();
	const [messageApi, contextHolder] = message.useMessage();

	const { trigger } = useUploadBundleMutation();

	const onFinish: FormProps<FieldType>["onFinish"] = async (values) => {
		try {
			await trigger({
				bundle: values.bundle[0].originFileObj as File,
				app_id: values.app_id,
				version_name: values.version_name,
				description: values.description ?? "",
			});
			navigate("/app");
		} catch (e) {
			messageApi.error(`${e}`);
		}
	};

	const onFinishFailed: FormProps<FieldType>["onFinishFailed"] = async () => {};

	return (
		<div style={{ padding: "24px 0" }}>
			{contextHolder}
			<h1>Create Bundle</h1>
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

					<Form.Item<FieldType> label="Description" name="description">
						<Input />
					</Form.Item>

					<Form.Item<FieldType>
						label="Bundle File"
						name="bundle"
						valuePropName="fileList"
						getValueFromEvent={(e) => (Array.isArray(e) ? e : e?.fileList)}
						rules={[{ required: true, message: "Please choose a bundle file!" }]}
					>
						<Upload
							accept="application/zip"
							multiple={false}
							maxCount={1}
							customRequest={({ onSuccess }) => onSuccess?.("")}
						>
							<Button icon={<UploadOutlined />}>Upload zip only</Button>
						</Upload>
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
