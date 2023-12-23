import logging
import httpx
import functools
from tenacity import retry, stop_after_attempt, wait_exponential_jitter, after_log
from typing import List, Callable, Awaitable
from typing_extensions import ParamSpec
from src.core.exceptions import HttpResponseError
from src.core.consts import (
    MAX_REQUEST_ATTEMPTS, LOG_HTTP_RESPONSES,
    INITIAL_RETRY_WAIT_TIME, MAX_RETRY_WAIT_TIME, RETRY_WAIT_EXP_BASE, RETRY_WAIT_JITTER
)

logger = logging.getLogger(__name__)
Params = ParamSpec("P")
HttpRequestCallableT = Callable[[Params], Awaitable[httpx.Response]]
CustomHttpRequestCallableT = Callable[[Params], Awaitable[dict]]


class StatusCodes:
    OK = 200
    CREATED = 201


def construct_response_details(response: httpx.Response) -> str:
    return f"{response.request.method} {response.request.url} [{response.status_code}]: {response.text}"


def validate_response(response: httpx.Response, expected_status_codes: List[int], log_response: bool):
    response_details = construct_response_details(response)

    if response.status_code in expected_status_codes:
        if log_response:
            logger.debug(response_details)
    else:
        raise HttpResponseError(response_details)


def http_request(expected_status_codes: List[int], max_attempts: int = MAX_REQUEST_ATTEMPTS,
                 log_response: bool = LOG_HTTP_RESPONSES):
    def decorator(function: HttpRequestCallableT) -> CustomHttpRequestCallableT:
        @functools.wraps(function)
        @retry(
            reraise=True,
            stop=stop_after_attempt(max_attempts),
            wait=wait_exponential_jitter(
                initial=INITIAL_RETRY_WAIT_TIME, max=MAX_RETRY_WAIT_TIME,
                exp_base=RETRY_WAIT_EXP_BASE, jitter=RETRY_WAIT_JITTER
            ),
            after=after_log(logger, logging.WARNING),
        )
        async def wrapper(*args: Params.args, **kwargs: Params.kwargs) -> dict:
            response = await function(*args, **kwargs)
            validate_response(response, expected_status_codes, log_response)
            return response.json()

        return wrapper

    return decorator
