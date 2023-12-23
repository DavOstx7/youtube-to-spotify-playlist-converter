import asyncio
from src.converter import PlaylistsConverter


async def main():
    converter = PlaylistsConverter()
    await converter.setup()
    await converter.run()


if __name__ == "__main__":
    asyncio.run(main())
