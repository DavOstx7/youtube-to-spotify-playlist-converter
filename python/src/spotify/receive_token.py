import uvicorn
import urllib.parse
from fastapi import FastAPI
from fastapi.responses import RedirectResponse
from src.core.utils import run_async
from src.spotify import api
from src.config.api_config import SPOTIFY_AUTHORIZATION_URL
from src.config.token_config import SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, SPOTIFY_REDIRECT_URI


async def manual():
    # step 1: run the function below, authorize against spotify, grab the code from the redirect uri query parameter.

    # api.authorize_via_browser(SPOTIFY_CLIENT_ID, SPOTIFY_REDIRECT_URI)

    code = "?"

    # step 2: run the functions below, only after you have set the code variable with the value from step 1.
    # p.s do not forget to comment out the function from step 1.

    # response = await api.request_access_token(SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, code, SPOTIFY_REDIRECT_URI)
    # print({"access_token": response["access_token"]})


def server():
    parsed_redirect_uri = urllib.parse.urlparse(SPOTIFY_REDIRECT_URI)

    app = FastAPI()

    @app.get("/")
    async def authorize():
        query_params = api.get_authorization_query_params(SPOTIFY_CLIENT_ID, SPOTIFY_REDIRECT_URI)
        return RedirectResponse(f"{SPOTIFY_AUTHORIZATION_URL}?{urllib.parse.urlencode(query_params)}")

    @app.get(parsed_redirect_uri.path)
    async def access_token(code: str):
        response = await api.request_access_token(SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, code, SPOTIFY_REDIRECT_URI)
        return {"access_token": response['access_token']}

    uvicorn.run(app, port=parsed_redirect_uri.port, host=parsed_redirect_uri.hostname)


if __name__ == "__main__":
    # run_async(manual())
    server()
