import logger from "../common/logger";
import api from "./api";
import { PlaylistQueryParams, PlaylistItemsPage } from "./models";
import { MAX_PAGE_RESULTS } from "../common/defaults";

export default class YouTubePlaylist {
    private queryParams: PlaylistQueryParams;
    private currentPage: PlaylistItemsPage | undefined;

    constructor(apiKey: string, playlistId: string, maxResults: number = MAX_PAGE_RESULTS) {
        this.queryParams = api.getPlaylistQueryParams(apiKey, playlistId, maxResults);
        this.currentPage = undefined;
    }

    get isInInitialState(): boolean {
        return this.currentPage === undefined;
    }

    get hasPrevPage(): boolean {
        if (this.isInInitialState) {
            return false;
        }
        return this.currentPage!.prevPageToken !== undefined;
    }

    get hasNextPage(): boolean {
        if (this.isInInitialState) {
            return false;
        }
        return this.currentPage!.nextPageToken !== undefined;
    }

    public async *walkPages(): AsyncIterable<PlaylistItemsPage> {
        yield await this.searchForPage();

        while (this.hasNextPage) {
            this.setNextPage();
            yield await this.searchForPage();
        }
    }

    public async searchForPage(): Promise<PlaylistItemsPage> {
        if (this.isInInitialState) {
            logger.debug(`Searching initial YouTube playlist page for playlist '${this.queryParams.playlistId}'`);
        } else {
            logger.debug(`Searching YouTube playlist page for playlist '${this.queryParams.playlistId}' with page token '${this.queryParams.pageToken}'`);
        }

        this.currentPage = await api.requestPlaylistPage(this.queryParams);
        return this.currentPage;
    }

    public setPrevPage() {
        this.queryParams.pageToken = this.currentPage.nextPageToken;
    }

    public setNextPage() {
        this.queryParams.pageToken = this.currentPage.nextPageToken;
    }
}
