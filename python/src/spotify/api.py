import httpx
import urllib.parse
import base64
import webbrowser
from typing import List
from src.core.http import http_request, StatusCodes
from src.spotify.models import AuthorizationQueryParams
from src.core.exceptions import ValidationError
from src.core.consts import REQUEST_TIMEOUT
from src.config.api_config import (
    SPOTIFY_MAX_TRACKS_PER_REQUEST, SPOTIFY_AUTHORIZATION_SCOPES,
    SPOTIFY_TOKEN_URL, SPOTIFY_USER_PROFILE_URL, SPOTIFY_PLAYLISTS_URL,
    SPOTIFY_SEARCH_URL, SPOTIFY_TRACKS_URL, SPOTIFY_AUTHORIZATION_URL
)

client = httpx.AsyncClient(timeout=httpx.Timeout(REQUEST_TIMEOUT))
MIN_POSITIVE_VALUE = 1


@http_request(expected_status_codes=[StatusCodes.OK], max_attempts=1)
async def request_access_token(client_id: str, client_secret: str, code: str, redirect_uri: str):
    headers = _get_request_access_token_headers(client_id, client_secret)
    form_body = {"grant_type": "authorization_code", "code": code, "redirect_uri": redirect_uri}
    return await client.post(SPOTIFY_TOKEN_URL, headers=headers, data=form_body)


@http_request(expected_status_codes=[StatusCodes.OK], max_attempts=1)
async def request_user_profile(access_token: str):
    headers = {"Authorization": f"Bearer {access_token}"}
    return await client.get(SPOTIFY_USER_PROFILE_URL, headers=headers)


@http_request(expected_status_codes=[StatusCodes.CREATED])
async def request_to_create_playlist(access_token: str, user_id: str, name: str, description: str, public: bool):
    url = SPOTIFY_PLAYLISTS_URL.format(user_id=user_id)
    headers = {"Content-Type": "application/json", "Authorization": f"Bearer {access_token}"}
    json_body = {"name": name, "description": description, "public": public}
    return await client.post(url, headers=headers, json=json_body)


@http_request(expected_status_codes=[StatusCodes.OK])
async def request_to_search_for_track(access_token: str, name: str, limit: int):
    query_params = {"q": name, "type": "track", "limit": limit}
    headers = {"Authorization": f"Bearer {access_token}"}
    return await client.get(SPOTIFY_SEARCH_URL, params=query_params, headers=headers)


@http_request(expected_status_codes=[StatusCodes.CREATED])
async def request_to_add_tracks(access_token: str, playlist_id: str, track_uris: List[str], position: int):
    _validate_track_uris_size(track_uris)
    url = SPOTIFY_TRACKS_URL.format(playlist_id=playlist_id)
    headers = {"Content-Type": "application/json", "Authorization": f"Bearer {access_token}"}
    json_body = {"uris": track_uris, "position": position}
    return await client.post(url, headers=headers, json=json_body)


def authorize_via_browser(client_id: str, redirect_uri: str):
    query_params = get_authorization_query_params(client_id, redirect_uri)
    webbrowser.open(f"{SPOTIFY_AUTHORIZATION_URL}?{urllib.parse.urlencode(query_params)}")


def get_authorization_query_params(client_id: str, redirect_uri: str) -> AuthorizationQueryParams:
    return {
        "client_id": client_id,
        "response_type": "code",
        "redirect_uri": redirect_uri,
        "scope": SPOTIFY_AUTHORIZATION_SCOPES
    }


def _get_request_access_token_headers(client_id: str, client_secret: str) -> dict:
    auth = f"{client_id}:{client_secret}"
    auth_64 = base64.urlsafe_b64encode(auth.encode()).decode()
    return {
        "Authorization": f"Basic {auth_64}",
        "Content-Type": "application/x-www-form-urlencoded"
    }


def _validate_track_uris_size(track_uris: List[str]):
    track_uris_size = len(track_uris)
    if not (MIN_POSITIVE_VALUE <= track_uris_size <= SPOTIFY_MAX_TRACKS_PER_REQUEST):
        valid_range = f"{MIN_POSITIVE_VALUE}-{SPOTIFY_MAX_TRACKS_PER_REQUEST}"
        raise ValidationError(f"The size of track uris ({track_uris_size}) is not in the valid range of {valid_range}")
