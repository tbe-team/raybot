from datetime import timedelta

from pydantic import Field

from .base import ConfigModel


class Credentials(ConfigModel):
    """Credentials for authentication"""

    username: str
    password: str


class AuthConfig(ConfigModel):
    """Auth configuration"""

    credentials: Credentials
    session_expiration_time: timedelta = Field(default=timedelta(seconds=10))
