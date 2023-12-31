from typing import List, Optional
from typing_extensions import TypedDict


class PlaylistQueryParams(TypedDict):
    key: str
    part: str
    playlistId: str
    maxResults: int
    pageToken: Optional[str]


class PlaylistItemSnippet(TypedDict, total=False):
    publishedAt: str
    channelId: str
    title: str
    description: str
    thumbnails: dict
    channelTitle: str
    videoOwnerChannelTitle: str
    videoOwnerChannelId: str
    playlistId: str
    position: int
    resourceId: dict
    contentDetails: dict
    status: dict


class PlaylistItem(TypedDict, total=False):
    id: str
    snippet: PlaylistItemSnippet


class PlaylistItemsPage(TypedDict, total=False):
    id: str
    nextPageToken: Optional[str]
    prevPageToken: Optional[str]
    pageInfo: dict
    items: List[PlaylistItem]
