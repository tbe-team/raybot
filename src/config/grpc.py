from pydantic import Field

from .base import ConfigModel


class GRPCConfig(ConfigModel):
    """GRPC configuration

    Attributes:
        port: The port to listen on
        max_workers: Controls the thread pool size for processing requests.
                    If all workers are busy, new request wait in a queue.
    """

    port: int = Field(default=50051)
    max_workers: int = Field(default=10)
