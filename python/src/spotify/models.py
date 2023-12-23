from typing_extensions import TypedDict


class AuthorizationQueryParams(TypedDict):
    client_id: str
    response_type: str
    redirect_uri: str
    scope: str


class PlaylistInfo(TypedDict):
    name: str
    description: str
    public: bool
