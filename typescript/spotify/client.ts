import logger from "../common/logger";
import api from "./api";
import { PlaylistInfo } from "./models";

type SpotifyPlaylistIdsT = string[];
type SpotifyTrackUrisT = string[];
const SINGLE_TRACK = 1;

export default class SpotifyClient {
    private accessToken: string;
    private userId: string;
    
    constructor(accessToken: string) {
        this.accessToken = accessToken;
        this.userId = undefined;
    }

    public async setUserId() {
        const response = await api.requestUserProfile(this.accessToken);
        this.userId = response.id;
    }

    public async createPlaylists(playlistsInfo: PlaylistInfo[]): Promise<SpotifyPlaylistIdsT> {
        const createPlaylistPromises = playlistsInfo.map((playlistInfo) => this.createPlaylist(playlistInfo));
        return await Promise.all(createPlaylistPromises);
    }

    public async searchForTrackUris(trackNames: string[]): Promise<SpotifyTrackUrisT> {
        logger.info(`Starting to search Spotify track uris for ${trackNames.length} track names...`);

        const searchForTrackUriPromsies = trackNames.map((trackName) => this.searchForTrackUri(trackName));
        let trackUris = await Promise.all(searchForTrackUriPromsies);
        trackUris = trackUris.filter((trackUri) => trackUri !== undefined);

        if (trackUris.length == 0) {
            logger.warning("Could not find a single Spotify track uri for the given track names");
        }

        return trackUris;
    }

    public async addTracksToPlaylist(playlistId: string, trackUris: string[], position: number = 0): Promise<string> {
        logger.info(`Adding ${trackUris.length} track uris to Spotify playlist '${playlistId}'`);
        const response = await api.requestToAddTracks(this.accessToken, playlistId, trackUris, position);
        logger.debug(`Snapshot id of the new added tracks for Spotify playlist '${playlistId}' is '${response.snapshot_id}'`);
        return response.snapshot_id;
    }

    private async createPlaylist(playlistInfo: PlaylistInfo): Promise<string> {
        const _playlistType = playlistInfo.public ? "public" : "private";
        logger.info(`Creating new ${_playlistType} Spotify playlist '${playlistInfo.name}'`);

        const response = await api.requestToCreatePlaylist(this.accessToken, this.userId, playlistInfo);
        return response.id;
    }

    private async searchForTrackUri(trackName: string): Promise<string | undefined> {
        const response = await api.requestToSearchForTracks(this.accessToken, trackName, SINGLE_TRACK);

        if ("tracks" in response) {
            const items = response.tracks.items;
            if (items.length > 0) {
                const trackUri = items[0]?.uri ?? undefined;

                if (trackUri) {
                    logger.debug(`Found Spotify track uri '${trackUri}' for '${trackName}'`);
                } else {
                    logger.warning(`Failed to find Spotify track uri for '${trackName}'`);
                }

                return trackUri;
            }   
        }

        return undefined;
    }
}
