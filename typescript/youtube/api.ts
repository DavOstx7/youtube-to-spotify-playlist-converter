import { URLSearchParams } from 'url';
import { HttpRequest, StatusCodes, fetchWithTimeout } from '../common/http';
import { PlaylistQueryParams } from './models';
import { ValidationError } from '../common/errors';
import apiConfig from '../../config/api_config.json';

const config = apiConfig.youtube;
const MIN_POSITIVE_VALUE = 1;


export class YoutubeAPI {
    @HttpRequest([StatusCodes.OK])
    public async requestPlaylistPage(queryParams: PlaylistQueryParams): Promise<any> {
        const {maxResults, ...rest} = queryParams;
        const searchParams = new URLSearchParams({...rest, maxResults: String(maxResults)});
        const url = `${config.urls.playlist_items}?${searchParams}`;

        return await fetchWithTimeout(url);
    }

    public getPlaylistQueryParams(apiKey: string, playlistId: string, maxResults: number): PlaylistQueryParams {
        this.validateMaxResultsValue(maxResults);
        return {"key": apiKey, "part": "snippet", "playlistId": playlistId, "maxResults": maxResults};
    }

    
    private validateMaxResultsValue(maxResults: number) {
        if (! this.isValidMaxResultsValue(maxResults)) {
            const validRange = `${MIN_POSITIVE_VALUE}-${config.max_items_per_request}`;
            throw new ValidationError(`The value of max results (${maxResults}) is not in the valid range of ${validRange}`);
        }
    }

    private isValidMaxResultsValue(maxResults: number): boolean {
        return maxResults >= MIN_POSITIVE_VALUE && maxResults <= config.max_items_per_request;
    }
}

const api = new YoutubeAPI();
export default api;
