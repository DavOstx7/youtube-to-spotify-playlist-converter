from typing import List
from src.core.utils import gather_async, configure_logging
from src.youtube.client import YoutubeClient
from src.spotify.client import SpotifyClient
from src.config.user_config import (
    LOGGING_LEVEL,
    YOUTUBE_API_KEY, YOUTUBE_PLAYLIST_IDS,
    SPOTIFY_ACCESS_TOKEN, SPOTIFY_NEW_PLAYLISTS, SPOTIFY_EXISTING_PLAYLIST_IDS
)

SpotifySnapshotIdsT = List[str]
SpotifySnapshotIdsGroupsT = List[SpotifySnapshotIdsT]


class PlaylistsConverter:
    def __init__(self):
        self._youtube_client = YoutubeClient(YOUTUBE_API_KEY)
        self._spotify_client = SpotifyClient(SPOTIFY_ACCESS_TOKEN)
        self._spotify_playlist_ids = []

    async def run(self) -> SpotifySnapshotIdsGroupsT:
        add_titles_to_playlists_tasks = [self._add_titles_to_playlists(titles_batch) async for titles_batch
                                         in self._youtube_client.walk_playlists_titles(YOUTUBE_PLAYLIST_IDS)]

        return await gather_async(*add_titles_to_playlists_tasks, filter_out_empty_results=True)

    async def setup(self):
        configure_logging(LOGGING_LEVEL)
        await self._spotify_client.set_user_id()
        await self._set_spotify_playlist_ids()

    async def _set_spotify_playlist_ids(self):
        self._spotify_playlist_ids.extend(SPOTIFY_EXISTING_PLAYLIST_IDS)

        if SPOTIFY_NEW_PLAYLISTS:
            new_playlist_ids = await self._spotify_client.create_playlists(SPOTIFY_NEW_PLAYLISTS)
            self._spotify_playlist_ids.extend(new_playlist_ids)

    async def _add_titles_to_playlists(self, titles_batch: List[str]) -> SpotifySnapshotIdsT:
        track_uris = await self._spotify_client.search_for_track_uris(titles_batch)

        add_tracks_to_playlist_tasks = [self._spotify_client.add_tracks_to_playlist(spotify_playlist_id, track_uris)
                                        for spotify_playlist_id in self._spotify_playlist_ids]

        return await gather_async(*add_tracks_to_playlist_tasks)
