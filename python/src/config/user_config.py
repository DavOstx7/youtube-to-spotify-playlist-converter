from src.core.utils import load_config_file

config = load_config_file("user_config.json")

LOGGING_LEVEL: str = config['logging']['level']
YOUTUBE_API_KEY: str = config['youtube']['api_key']
YOUTUBE_PLAYLIST_IDS: list = config['youtube']['playlist_ids']
SPOTIFY_ACCESS_TOKEN: str = config['spotify']['access_token']
SPOTIFY_NEW_PLAYLISTS: list = config['spotify']['new_playlists']
SPOTIFY_EXISTING_PLAYLIST_IDS: list = config['spotify']['existing_playlist_ids']
