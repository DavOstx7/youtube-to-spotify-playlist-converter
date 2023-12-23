from src.core.utils import load_config_file

config = load_config_file("api_config.json")

YOUTUBE_PLAYLIST_ITEMS_URL: str = config['youtube']['urls']['playlist_items']
YOUTUBE_MAX_ITEMS_PER_REQUEST: int = config['youtube']['max_items_per_request']
SPOTIFY_USER_PROFILE_URL: str = config['spotify']['urls']['user_profile']
SPOTIFY_SEARCH_URL: str = config['spotify']['urls']['search']
SPOTIFY_TOKEN_URL: str = config['spotify']['urls']['token']
SPOTIFY_AUTHORIZATION_URL: str = config['spotify']['urls']['authorization']
SPOTIFY_PLAYLISTS_URL: str = config['spotify']['urls']['playlists']
SPOTIFY_TRACKS_URL: str = config['spotify']['urls']['tracks']
SPOTIFY_AUTHORIZATION_SCOPES: str = config['spotify']['authorization_scopes']
SPOTIFY_MAX_TRACKS_PER_REQUEST: int = config['spotify']['max_tracks_per_request']
