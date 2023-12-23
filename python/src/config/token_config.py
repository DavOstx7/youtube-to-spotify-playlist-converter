from src.core.utils import load_config_file

config = load_config_file("token_config.json")

SPOTIFY_CLIENT_ID: str = config['spotify']['client_id']
SPOTIFY_CLIENT_SECRET: str = config['spotify']['client_secret']
SPOTIFY_REDIRECT_URI: str = config['spotify']['redirect_uri']
