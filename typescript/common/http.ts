import pRetry, { FailedAttemptError } from "p-retry";
import logger from "./logger";
import { ResponseError, TimeoutError } from "./errors";
import { 
    REQUEST_TIMEOUT_MILLISECONDS, LOG_HTTP_RESPONSES, MAX_REQUEST_RETRIES,
    MIN_RETRY_TIMEOUT_MILLISECONDS, MAX_RETRY_TIMEOUT_MILLISECONDS,
    RETRY_EXP_FACTOR_MILLISECONDS, USE_RETRY_JITTER

} from "./defaults";

interface TimedRequestInit extends RequestInit {
    timeout?: number;
} 

export const StatusCodes = {
    OK: 200,
    CREATED: 201,
} as const;
  
function constructResponseDetails(response: Response, responseJson: any): string {
    const responseText = JSON.stringify(responseJson, null, 2);
    return `${response.url} [${response.status}]: ${responseText}`;
}

async function validateResponse(response: Response, expectedStatusCodes: number[], logResponse: boolean): Promise<any> {
    const responseJson = await response.json();
    const responseDetails = constructResponseDetails(response, responseJson);

    if (expectedStatusCodes.includes(response.status)) {
        if (logResponse) {
            logger.debug(responseDetails);
        }
        return responseJson;
    }

    throw new ResponseError(responseDetails);
}

export function HttpRequest(expectedStatusCodes: number[], maxRetries: number = MAX_REQUEST_RETRIES, logResponse: boolean = LOG_HTTP_RESPONSES) {
    return function (target: any, propertyKey: string, descriptor: PropertyDescriptor) { 
        const originalMethod = descriptor.value;

        descriptor.value = async function(...args: any[]): Promise<any> {
            return await pRetry(async () => {
                const response = await originalMethod.apply(this, args);
                return await validateResponse(response, expectedStatusCodes, logResponse);
            },
            {
                retries: maxRetries,
                onFailedAttempt: (error: FailedAttemptError) => {
                    logger.warn(`Attempt ${error.attemptNumber} failed. There are ${error.retriesLeft} retries left`);
                },
                minTimeout: MIN_RETRY_TIMEOUT_MILLISECONDS,
                maxTimeout: MAX_RETRY_TIMEOUT_MILLISECONDS,
                factor: RETRY_EXP_FACTOR_MILLISECONDS,
                randomize: USE_RETRY_JITTER,    
            });
        }

        return descriptor;
    };
}

export async function fetchWithTimeout(url: string, options?: TimedRequestInit): Promise<Response> {
    const timeout = options?.timeout ?? REQUEST_TIMEOUT_MILLISECONDS;

    const controller = new AbortController();
    const timeoutError = new TimeoutError(`Fetch timed out after ${timeout} milliseconds`);
    const timeoutId = setTimeout(() => controller.abort(timeoutError), timeout);

    const response = await fetch(url, {
        ...options,
        signal: controller.signal  
    });
    
    clearTimeout(timeoutId);
    return response;
}
