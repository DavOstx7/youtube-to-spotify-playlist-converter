import asyncio
import sys
import os
import json
import logging.config
from typing import TypeVar, Awaitable

_FILE_DIR = os.path.dirname(__file__)
_CONFIG_DIR = os.path.join(_FILE_DIR, "../../../config/")

T = TypeVar("T")


def run_async(main: Awaitable[T]) -> T:
    if sys.version_info >= (3, 7):
        return asyncio.run(main)
    else:
        loop = asyncio.get_event_loop()
        try:
            return loop.run_until_complete(main)
        finally:
            loop.close()


async def gather_async(*coros_or_futures: Awaitable, return_exceptions=False,
                       filter_out_empty_results: bool = False) -> list:
    results = await asyncio.gather(*coros_or_futures, return_exceptions=return_exceptions)

    if filter_out_empty_results:
        return list(filter(lambda result: bool(result), results))
    else:
        return results


def load_config_file(filename: str):
    config_file = os.path.join(_CONFIG_DIR, filename)
    with open(config_file, "r") as file:
        return json.load(file)


def configure_logging(level: str):
    logging.config.dictConfig({
        "version": 1,
        "disable_existing_loggers": True,
        "formatters": {
            "standard": {
                "format": "%(name)s : %(asctime)s : %(levelname)s : %(message)s",
                "datefmt": "%Y-%m-%dT%H:%M:%S"
            },
        },
        "handlers": {
            "stream": {
                "level": "DEBUG",
                "formatter": "standard",
                "class": "logging.StreamHandler",
                "stream": "ext://sys.stdout",
            },
            "null": {
                "class": "logging.NullHandler",
            },
        },
        "root": {
            "handlers": ["null"]
        },
        "loggers": {
            "src": {
                "handlers": ["stream"],
                "level": level,
                "propagate": False
            }
        }
    })
