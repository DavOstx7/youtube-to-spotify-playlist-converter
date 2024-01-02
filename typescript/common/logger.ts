import winston from "winston";
const { combine, timestamp, printf, colorize } = winston.format;

const logger = winston.createLogger({
  level: "info",
  format: combine(
    colorize({ all: true }),
    timestamp({format: "YYYY-MM-DD hh:mm:ss"}),
    printf((entry: winston.LogEntry) => `[${entry.timestamp}] : ${entry.level} : ${entry.message}`)
  ),
  transports: [new winston.transports.Console()],
});

export default logger;
