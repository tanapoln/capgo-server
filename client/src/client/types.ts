

export type Platform = "ios" | "android"

export type UploadBundleRequest = {
    bundle: File;
    version_name: string;
    description: string;
}

export type BundleResponse = {
    id: string;
    version_name: string;
    description: string;
    crc_checksum: string;
    public_download_url: string;
    created_at: string;
}

export type ListAllBundlesResponse = {
    data: BundleResponse[];
}

export type BundleModifiedResponse = {
    message: string;
    bundle: BundleResponse;
}

export type ReleaseResponse = {
    id: string;
    platform: Platform;
    version_name: string;
    version_code: string;
    release_date: string | null | undefined;
    builtin_bundle_id: string
    active_bundle_id: string | null | undefined;
    created_at: string
}

export type ListAllReleasesResponse = {
    data: ReleaseResponse[]
}

export type ReleaseModifiedResponse = {
    message: string;
    release: ReleaseResponse;
}

export type SetReleaseActiveBundleRequest = {
    release_id: string;
    bundle_id: string;
}

export type CreateReleaseRequest = {
    platform: Platform;
    app_id: string;
    version_name: string;
    version_code: string;
    builtin_bundle_id: string;
}

export type UpdateReleaseRequest = {
    release_id: string;
    release_date: string | undefined;
}
