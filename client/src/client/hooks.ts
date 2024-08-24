import useSWR, { mutate } from "swr";
import {
	createRelease,
	deleteRelease,
	listBundles,
	listReleases,
	setReleaseActiveBundle,
	updateRelease,
	uploadBundle,
} from "./api";
import type {
	CreateReleaseRequest,
	DeleteReleaseRequest,
	SetReleaseActiveBundleRequest,
	UpdateReleaseRequest,
	UploadBundleRequest,
} from "./types";

export function useReleases() {
	return useSWR("releases", listReleases);
}

export function useRelease(id: string) {
	const { data: releases } = useReleases();
	return useSWR(releases ? `releases/${id}` : null, async () => {
		const rec = releases!.data.find((release) => release.id === id);
		if (rec === undefined) {
			throw new Error(`Release ${id} not found`);
		}
		return rec;
	});
}

export function useBundles() {
	return useSWR("bundles", listBundles);
}

export function useBundle(id: string) {
	const { data: bundles } = useBundles();
	return useSWR(bundles ? `bundles/${id}` : null, async () => {
		const rec = bundles!.data.find((bundle) => bundle.id === id);
		if (rec === undefined) {
			throw new Error(`Bundle ${id} not found`);
		}
		return rec;
	});
}

export function useUploadBundleMutation() {
	const trigger = async (req: UploadBundleRequest) => {
		await uploadBundle(req);
		mutate("bundles");
	};
	return { trigger };
}

export function useCreateReleaseMutation() {
	const trigger = async (req: CreateReleaseRequest) => {
		await createRelease(req);
		mutate("releases");
	};
	return { trigger };
}

export function useUpdateReleaseMutation() {
	const trigger = async (req: UpdateReleaseRequest) => {
		await updateRelease(req);
		mutate("releases");
	};
	return { trigger };
}

export function useSetReleaseActiveBundleMutation() {
	const trigger = async (req: SetReleaseActiveBundleRequest) => {
		await setReleaseActiveBundle(req);
		mutate("releases");
	};
	return { trigger };
}

export function useDeleteReleaseMutation() {
	const trigger = async (req: DeleteReleaseRequest) => {
		await deleteRelease(req);
		mutate("releases");
	};
	return { trigger };
}
