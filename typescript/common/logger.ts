export enum LoggingLevel {
    DEBUG = 10,
    INFO = 20,
    WARNING = 30,
    ERROR = 40,
    CRITICAL = 50,
}

const LOGGING_LEVEL_TO_STRING = {
    10: "DEBUG",
    20: "INFO",
    30: "WARNING",
    40: "ERROR",
    50: "CRITICAL"
}

class Logger {
    private level: LoggingLevel;
    private includeTimestamp: boolean;

    constructor (level: LoggingLevel = LoggingLevel.INFO, includeTimestamp: boolean = true) {
        this.level = level;
        this.includeTimestamp = includeTimestamp;
    };

    public setLevel(level: number) {
        this.level = level;
    }

    public log(level: LoggingLevel, message: string) {
        if (level < this.level) {
            return;
        }

        let logMessage = this.formatLogMessage(level, message);
        console.log(logMessage);
    }
    
    public debug(message: string) {
        this.log(LoggingLevel.DEBUG, message);
    }

    public info(message: string) {
        this.log(LoggingLevel.INFO, message);
    }

    public warning(message: string) {
        this.log(LoggingLevel.WARNING, message);
    }

    public error(message: string) {
        this.log(LoggingLevel.ERROR, message);
    }

    public critical(message: string) {
        this.log(LoggingLevel.CRITICAL, message);
    }

    private formatLogMessage(level: LoggingLevel, message: string): string {
        let logMessage = "";

        if (this.includeTimestamp) {
            logMessage += `${this.getCurrentDatetime()} : `;
        }

        logMessage += `${LOGGING_LEVEL_TO_STRING[level]} : `;
        logMessage += message;
        return logMessage;
    }

    private getCurrentDatetime(): string {
        const date = new Date();
        const [datetime, milliseconds] =  date.toISOString().split('.');
        return datetime;
    }
}

const logger = new Logger();
export default logger;
