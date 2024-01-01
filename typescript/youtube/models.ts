export interface PlaylistQueryParams {
    key: string;
    part: string;
    playlistId: string;
    maxResults: number;
    pageToken?: string;
}

export interface PlaylistItemSnippet {
    publishedAt: string;
    channelId: string;
    title: string;
    description: string;
    thumbnails: Record<string, string | number>;
    channelTitle: string;
    videoOwnerChannelTitle: string;
    videoOwnerChannelId: string;
    playlistId: string;
    position: number;
    resourceId: Record<string, string>;
    contentDetails: Record<string, string>;
    status: Record<string, string>;

    [key: string]: unknown;
}

export interface PlaylistItem {
    id: string;
    snippet:  PlaylistItemSnippet;

    [key: string]: unknown;
}

export interface PlaylistItemsPage {
    id: string;
    nextPageToken?: string;
    prevPageToken?: string;
    pageInfo: Record<string, number>;
    items: PlaylistItem[];

    [key: string]: unknown;
}
