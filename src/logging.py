import logging as std_logging

from src.config.log import LogConfig, LogLevel


def setup_logging(config: LogConfig) -> std_logging.Logger:
    """Setup logging with the given config"""
    level_map = {
        LogLevel.DEBUG: std_logging.DEBUG,
        LogLevel.INFO: std_logging.INFO,
        LogLevel.WARNING: std_logging.WARNING,
        LogLevel.ERROR: std_logging.ERROR,
    }
    level = level_map.get(config.level, std_logging.INFO)

    console_handler = std_logging.StreamHandler()
    console_handler.setLevel(level)
    if config.colorize:
        console_handler.setFormatter(ColorizedFormatter())
    else:
        console_handler.setFormatter(PlainFormatter())
    std_logging.basicConfig(level=level, handlers=[console_handler])

    return std_logging.getLogger()


class PlainFormatter(std_logging.Formatter):
    """Plain formatter without color."""

    def __init__(self, fmt: str | None = None) -> None:
        super().__init__()
        # Loguru style without color
        self.fmt = fmt or (
            "%(asctime)s | "
            "%(levelname)-8s | "
            "%(name)s:%(funcName)s:%(lineno)d"
            " - %(message)s"
        )

    def format(self, record: std_logging.LogRecord) -> str:
        formatter = std_logging.Formatter(self.fmt)
        return formatter.format(record)


class ColorizedFormatter(std_logging.Formatter):
    """std_logging colored formatter, adapted from https://stackoverflow.com/a/56944256/3638629"""

    bold = "\x1b[1m"
    grey = f"{bold}\x1b[38;21m"
    blue = f"{bold}\x1b[38;5;39m"
    yellow = f"{bold}\x1b[38;5;229m"
    red = f"{bold}\x1b[38;5;203m"
    bold_red = f"{bold}\x1b[31;1m"

    def __init__(self, fmt: str | None = None) -> None:
        super().__init__()
        # Loguru color style
        self.fmt = fmt or (
            "\x1b[32m%(asctime)s\x1b[0m | "
            "%(levelname)-8s\x1b[0m | "
            "\x1b[36m%(name)s\x1b[0m:\x1b[36m%(funcName)s\x1b[0m:\x1b[36m%(lineno)d\x1b[0m"
            " - %(message)s\x1b[0m"
        )
        self.FORMATS = {
            std_logging.DEBUG: self.blue,
            std_logging.INFO: self.grey,
            std_logging.WARNING: self.yellow,
            std_logging.ERROR: self.red,
            std_logging.CRITICAL: self.bold_red,
        }

    def format(self, record: std_logging.LogRecord) -> str:
        log_fmt = self.FORMATS.get(record.levelno, "")
        formatter = std_logging.Formatter(self.fmt)
        record.levelname = log_fmt + record.levelname
        record.msg = log_fmt + str(record.msg)
        return formatter.format(record)
