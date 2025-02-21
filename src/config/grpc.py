from pydantic import Field

from .base import ConfigModel


class GRPCConfig(ConfigModel):
    """GRPC configuration"""

    port: int = Field(default=50051)
