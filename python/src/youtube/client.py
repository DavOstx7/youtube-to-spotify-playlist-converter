import logging
from typing import List, AsyncGenerator
from src.youtube.playlist import YoutubePlaylist
from src.core.consts import MAX_TITLES_BATCH_SIZE

logger = logging.getLogger(__name__)
AsyncYoutubeTitlesIteratorT = AsyncGenerator[List[str], None]


class YoutubeClient:
    def __init__(self, api_key: str):
        self._api_key = api_key

    async def walk_playlists_titles(self, playlist_ids: List[str],
                                    max_batch_size: int = MAX_TITLES_BATCH_SIZE) -> AsyncYoutubeTitlesIteratorT:
        titles_batch = []

        for playlist_id in playlist_ids:
            playlist = YoutubePlaylist(self._api_key, playlist_id)

            logger.info(f"Starting to search for YouTube video titles inside playlist '{playlist_id}'...")
            async for page in playlist.walk_pages():
                for item in page['items']:
                    logger.debug(f"Found YouTube video title '{item['snippet']['title']}'")
                    titles_batch.append(item['snippet']['title'])

                    if len(titles_batch) >= max_batch_size:
                        yield titles_batch
                        titles_batch = []

        if titles_batch:
            yield titles_batch
