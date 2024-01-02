import logger from "./common/logger";
import YoutubeClient from "./youtube/client";
import SpotifyClient from "./spotify/client";
import userConfig from '../config/user_config.json';

type SpotifySnapshotIdsT = string[];
type SpotifySnapshotIdsGroupsT = SpotifySnapshotIdsT[];

export default class PlaylistsConverter {
    private youtubeClient: YoutubeClient;
    private spotifyClient: SpotifyClient;
    private spotifyPlaylistIds: string[];

    constructor() {
        this.youtubeClient = new YoutubeClient(userConfig.youtube.api_key);
        this.spotifyClient = new SpotifyClient(userConfig.spotify.access_token);
        this.spotifyPlaylistIds = [];
    }

    public async run(): Promise<SpotifySnapshotIdsGroupsT> {
        const addTitlesToPlaylistsPromsises = [];
        for await (const titlesBatch of this.youtubeClient.walkPlaylistsTitles(userConfig.youtube.playlist_ids)) {
            addTitlesToPlaylistsPromsises.push(this.addTitlesToPlaylists(titlesBatch));
        }

        let spotifySnapshotIdsGroups = await Promise.all(addTitlesToPlaylistsPromsises);
        return spotifySnapshotIdsGroups.filter((spotifySnapshotIdsGroup) => spotifySnapshotIdsGroup.length > 0);

    }

    public async setup() {
        logger.level = userConfig.logging.level.toLowerCase();
        await this.spotifyClient.setUserId();
        await this.setSpotifyPlaylistIds();
    }

    private async setSpotifyPlaylistIds() {
        this.spotifyPlaylistIds = this.spotifyPlaylistIds.concat(userConfig.spotify.existing_playlist_ids);

        if (userConfig.spotify.new_playlists.length > 0) {
            const newPlaylistIds = await this.spotifyClient.createPlaylists(userConfig.spotify.new_playlists);
            this.spotifyPlaylistIds = this.spotifyPlaylistIds.concat(newPlaylistIds);
        }
    }

    private async addTitlesToPlaylists(titlesBatch: string[]): Promise<SpotifySnapshotIdsT> {
        const trackUris = await this.spotifyClient.searchForTrackUris(titlesBatch);

        const addTracksToPlaylistPromises = this.spotifyPlaylistIds.map(
            (spotifyPlaylistId) => this.spotifyClient.addTracksToPlaylist(spotifyPlaylistId, trackUris)
        );

        return await Promise.all(addTracksToPlaylistPromises);
    }
}