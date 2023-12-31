import logging
from typing import List, Optional
from src.core.utils import gather_async
from src.spotify import api
from src.spotify.models import PlaylistInfo

logger = logging.getLogger(__name__)
SpotifyPlaylistIdsT = List[str]
SpotifyTrackUrisT = List[str]
SINGLE_TRACK = 1


class SpotifyClient:
    def __init__(self, access_token: str):
        self._access_token = access_token
        self._user_id: Optional[str] = None

    async def set_user_id(self):
        response = await api.request_user_profile(self._access_token)
        self._user_id = response['id']

    async def create_playlists(self, playlists_info: List[PlaylistInfo]) -> SpotifyPlaylistIdsT:
        create_playlist_tasks = [self._create_playlist(playlist_info) for playlist_info in playlists_info]
        return await gather_async(*create_playlist_tasks)

    async def search_for_track_uris(self, track_names: List[str]) -> SpotifyTrackUrisT:
        logger.info(f"Starting to search Spotify track uris for {len(track_names)} track names...")

        search_for_track_uri_tasks = [self._search_for_track_uri(track_name) for track_name in track_names]
        track_uris = await gather_async(*search_for_track_uri_tasks, filter_out_empty_results=True)

        if not track_uris:
            logger.warning("Could not find a single Spotify track uri for the given track names")

        return track_uris

    async def add_tracks_to_playlist(self, playlist_id: str, track_uris: List[str], position: int = 0) -> str:
        logger.info(f"Adding {len(track_uris)} track uris to Spotify playlist '{playlist_id}'")
        response = await api.request_to_add_tracks(self._access_token, playlist_id, track_uris, position)
        logger.debug(
            f"Snapshot id of the new added tracks for Spotify playlist '{playlist_id}' is '{response['snapshot_id']}'"
        )
        return response['snapshot_id']

    async def _create_playlist(self, playlist_info: PlaylistInfo) -> str:
        _playlist_type = "public" if playlist_info['public'] else "private"
        logger.info(f"Creating new {_playlist_type} Spotify playlist '{playlist_info['name']}'")

        response = await api.request_to_create_playlist(self._access_token, self._user_id, **playlist_info)
        return response['id']

    async def _search_for_track_uri(self, track_name: str) -> Optional[str]:
        response = await api.request_to_search_for_track(self._access_token, track_name, SINGLE_TRACK)

        if 'tracks' in response:
            items = response['tracks']['items']
            if items:
                track_uri = items[0].get('uri')

                if track_uri:
                    logger.debug(f"Found Spotify track uri '{track_uri}' for '{track_name}'")
                else:
                    logger.warning(f"Failed to find Spotify track uri for '{track_name}'")

                return track_uri
