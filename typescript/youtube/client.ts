import logger from "../common/logger";
import YouTubePlaylist from "./playlist";
import { MAX_TITLES_BATCH_SIZE } from "../common/defaults";

export default class YoutubeClient {
    private apiKey: string;

    constructor(apiKey: string) {
        this.apiKey = apiKey;
    }

    public async *walkPlaylistsTitles(playlistIds: string[], maxBatchSize: number = MAX_TITLES_BATCH_SIZE): AsyncIterable<string[]> {
        let titlesBatch  = [];
        
        for (const playlistId of playlistIds) {
            const playlist = new YouTubePlaylist(this.apiKey, playlistId);

            logger.info(`Starting to search for YouTube video titles inside playlist '${playlistId}'...`);
            for await (const page of playlist.walkPages()) {
                for (const item of page.items) {
                    logger.debug(`Found YouTube video title '${item.snippet.title}'`);
                    titlesBatch.push(item.snippet.title);
                    
                    if (titlesBatch.length >= maxBatchSize){
                        yield titlesBatch;
                        titlesBatch = [];
                    }
                }
            }
        }

        if (titlesBatch.length > 0){
            yield titlesBatch;
        }
    }
}
