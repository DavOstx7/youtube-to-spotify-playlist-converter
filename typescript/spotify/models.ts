export interface AuthorizationQueryParams {
    client_id: string;
    response_type: string;
    redirect_uri: string;
    scope: string;

    [key: string]: string;
}

export interface PlaylistInfo {
    name: string;
    description: string;
    public: boolean;
}

export interface RequestAccessTokenHeaders {
    "Authorization": string;
    "Content-Type": string;

    [key: string]: string;
}
