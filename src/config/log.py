from enum import StrEnum

from pydantic import Field

from .base import ConfigModel


class LogLevel(StrEnum):
    DEBUG = "DEBUG"
    INFO = "INFO"
    WARNING = "WARNING"
    ERROR = "ERROR"


class LogConfig(ConfigModel):
    """Logging configuration"""

    level: LogLevel = Field(default=LogLevel.INFO)
    colorize: bool = Field(default=True)
