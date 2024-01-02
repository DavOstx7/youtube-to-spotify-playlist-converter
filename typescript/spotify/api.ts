import { URLSearchParams } from 'url';
import { HttpRequest, StatusCodes, fetchWithTimeout } from '../common/http';
import { AuthorizationQueryParams, PlaylistInfo } from './models';
import { ValidationError } from '../common/errors';
import apiConfig from '../../config/api_config.json';

const config = apiConfig.spotify;
const NO_RETRIES = 0;
const MIN_POSITIVE_VALUE = 1;

export class SpotifyAPI {
    @HttpRequest([StatusCodes.OK], NO_RETRIES)
    public async requestAccessToken(clientId: string, clientSecret: string, code: string, redirectUri: string): Promise<any> {
        const headers = this.getRequestAccessTokenHeaders(clientId, clientSecret);
        const formBody = new URLSearchParams({"grant_type": "authorization_code", "code": code, "redirect_uri": redirectUri}).toString();

        return await fetchWithTimeout(config.urls.token, {
            method: "POST",
            headers: headers,
            body: formBody
        });
    }

    @HttpRequest([StatusCodes.OK], NO_RETRIES)
    public async requestUserProfile(accessToken: string): Promise<any> {
        const headers = {"Authorization": `Bearer ${accessToken}`};

        return await fetchWithTimeout(config.urls.user_profile, {
            method: "GET",
            headers: headers
        });
    }

    @HttpRequest([StatusCodes.CREATED])
    public async requestToCreatePlaylist(accessToken: string, userId: string, {name, description, public: isPublic}: PlaylistInfo): Promise<any> {
        const url = config.urls.playlists.replace("{user_id}", userId);
        const headers = {"Content-Type": "application/json", "Authorization": `Bearer ${accessToken}`};
        const jsonBody = JSON.stringify({"name": name, "description": description, "public": isPublic});

        return await fetchWithTimeout(url, {
            method: "POST",
            headers: headers,
            body: jsonBody
        });
    }

    @HttpRequest([StatusCodes.OK])
    public async requestToSearchForTracks(accessToken: string, name: string, limit: number): Promise<any> {
        const searchParams = new URLSearchParams({"q": name, "type": "track", "limit": String(limit)});
        const url = `${config.urls.search}?${searchParams}`;
        const headers = {"Authorization": `Bearer ${accessToken}`};

        return await fetchWithTimeout(url, {
            method: "GET",
            headers: headers
        });
    }

    @HttpRequest([StatusCodes.CREATED])
    public async requestToAddTracks(accessToken: string, playlistId: string, trackUris: string[], position: number): Promise<any> {
        this.validateTrackUrisSize(trackUris);
        
        const url = config.urls.tracks.replace("{playlist_id}", playlistId);
        const headers = {"Content-Type": "application/json", "Authorization": `Bearer ${accessToken}`};
        const jsonBody = JSON.stringify({"uris": trackUris, "position": position});

        return await fetchWithTimeout(url, {
            method: "POST",
            headers: headers,
            body: jsonBody
        });
    }

    public getAuthorizationQueryParams(clientId: string, redirectUri: string): AuthorizationQueryParams {
        return {
            client_id: clientId,
            response_type: "code",
            redirect_uri: redirectUri,
            scope: apiConfig.spotify.authorization_scopes
        };
    }

    private getRequestAccessTokenHeaders(clientId: string, clientSecret: string): Record<string, string> {
        const auth = `${clientId}:${clientSecret}`;
        const auth64 = Buffer.from(auth).toString("base64");

        return {
            "Authorization": `Basic ${auth64}`,
            "Content-Type": "application/x-www-form-urlencoded"
        };
    }

    private validateTrackUrisSize(trackUris: string[]) {
        const trackUrisSize = trackUris.length;
        if (! this.isValidTrackUrisSize(trackUrisSize)) {
            const validRange = `${MIN_POSITIVE_VALUE}-${config.max_tracks_per_request}`;
            throw new ValidationError(`The size of track uris (${trackUrisSize}) is not in the valid range of ${validRange}`);
        }
    }

    private isValidTrackUrisSize(trackUrisSize: number): boolean {
        return trackUrisSize >= MIN_POSITIVE_VALUE && trackUrisSize <= config.max_tracks_per_request;
    }
}

const api = new SpotifyAPI();
export default api;
