import httpx
from src.core.http import http_request, StatusCodes
from src.youtube.models import PlaylistQueryParams
from src.core.exceptions import ValidationError
from src.core.consts import REQUEST_TIMEOUT
from src.config.api_config import YOUTUBE_MAX_ITEMS_PER_REQUEST, YOUTUBE_PLAYLIST_ITEMS_URL

client = httpx.AsyncClient(timeout=httpx.Timeout(REQUEST_TIMEOUT))
MIN_POSITIVE_VALUE = 1


@http_request(expected_status_codes=[StatusCodes.OK])
async def request_playlist_page(query_params: PlaylistQueryParams):
    return await client.get(YOUTUBE_PLAYLIST_ITEMS_URL, params=query_params)


def get_playlist_query_params(api_key: str, playlist_id: str, max_results: int) -> PlaylistQueryParams:
    _validate_max_results_value(max_results)
    return {"key": api_key, "part": "snippet", "playlistId": playlist_id, "maxResults": max_results}


def _validate_max_results_value(max_results: int):
    if not (MIN_POSITIVE_VALUE <= max_results <= YOUTUBE_MAX_ITEMS_PER_REQUEST):
        valid_range = f"{MIN_POSITIVE_VALUE}-{YOUTUBE_MAX_ITEMS_PER_REQUEST}"
        raise ValidationError(f"The value of max results ({max_results}) is not in the valid range of {valid_range}")
