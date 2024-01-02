import PlaylistsConverter from "./converter";

async function main() {
    const playlistsConverter = new PlaylistsConverter();
    await playlistsConverter.setup();
    await playlistsConverter.run();
}

main();
