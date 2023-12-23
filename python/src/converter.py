import asyncio
from typing import List
from src.core.utils import configure_logging
from src.youtube.client import YoutubeClient
from src.spotify.client import SpotifyClient
from src.config.user_config import (
    LOGGING_LEVEL,
    YOUTUBE_API_KEY, YOUTUBE_PLAYLIST_IDS,
    SPOTIFY_ACCESS_TOKEN, SPOTIFY_NEW_PLAYLISTS, SPOTIFY_EXISTING_PLAYLIST_IDS
)

SpotifySnapshotIdsT = List[str]


class PlaylistsConverter:
    def __init__(self):
        self._youtube_client = YoutubeClient(YOUTUBE_API_KEY)
        self._spotify_client = SpotifyClient(SPOTIFY_ACCESS_TOKEN)
        self._spotify_playlist_ids = []

    async def run(self) -> List[SpotifySnapshotIdsT]:
        add_titles_to_playlists_tasks = []
        async for titles_batch in self._youtube_client.walk_playlists_titles(YOUTUBE_PLAYLIST_IDS):
            add_titles_to_playlists_tasks.append(self._add_tracks_to_spotify_playlists(titles_batch))

        await asyncio.gather(*add_titles_to_playlists_tasks)

    async def setup(self):
        configure_logging(LOGGING_LEVEL)
        await self._spotify_client.set_user_id()
        await self._set_spotify_playlist_ids()

    async def _set_spotify_playlist_ids(self):
        self._spotify_playlist_ids.extend(SPOTIFY_EXISTING_PLAYLIST_IDS)
        if SPOTIFY_NEW_PLAYLISTS:
            new_playlist_ids = await self._spotify_client.create_playlists(SPOTIFY_NEW_PLAYLISTS)
            self._spotify_playlist_ids.extend(new_playlist_ids)

    async def _add_tracks_to_spotify_playlists(self, track_names: List[str]) -> SpotifySnapshotIdsT:
        track_uris = await self._spotify_client.search_for_track_uris(track_names)

        if not track_uris:
            return

        add_tracks_to_playlist_tasks = []
        for spotify_playlist_id in self._spotify_playlist_ids:
            add_tracks_to_playlist_tasks.append(
                self._spotify_client.add_tracks_to_playlist(spotify_playlist_id, track_uris)
            )

        return await asyncio.gather(*add_tracks_to_playlist_tasks)
