import logging
from typing import AsyncGenerator
from src.youtube import api
from src.youtube.models import PlaylistItemsPage
from src.core.consts import MAX_PAGE_RESULTS

logger = logging.getLogger(__name__)
AsyncYoutubePlaylistItemsPageIteratorT = AsyncGenerator[PlaylistItemsPage, None]


class YoutubePlaylist:
    def __init__(self, api_key: str, playlist_id: str, max_results: int = MAX_PAGE_RESULTS):
        self._query_params = api.get_playlist_query_params(api_key, playlist_id, max_results)
        self._current_page: PlaylistItemsPage = None

    @property
    def is_in_initial_state(self) -> bool:
        return self._current_page is None

    @property
    def has_prev_page(self) -> bool:
        if self.is_in_initial_state:
            return False
        return 'prevPageToken' in self._current_page

    @property
    def has_next_page(self) -> bool:
        if self.is_in_initial_state:
            return False
        return 'nextPageToken' in self._current_page

    async def walk_pages(self) -> AsyncYoutubePlaylistItemsPageIteratorT:
        yield await self.search_for_page()

        while self.has_next_page:
            self.set_next_page()
            yield await self.search_for_page()

    async def search_for_page(self) -> PlaylistItemsPage:
        if self.is_in_initial_state:
            logger.debug(f"Searching initial YouTube playlist page for playlist '{self._query_params['playlistId']}'")
        else:
            logger.debug(f"Searching YouTube playlist page for playlist '{self._query_params['playlistId']}' "
                         f"with page token '{self._query_params['pageToken']}'")

        self._current_page = await api.request_playlist_page(self._query_params)
        return self._current_page

    def set_prev_page(self):
        self._query_params['pageToken'] = self._current_page['prevPageToken']

    def set_next_page(self):
        self._query_params['pageToken'] = self._current_page['nextPageToken']
