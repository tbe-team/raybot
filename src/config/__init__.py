from typing import Any, cast

import yaml
from pydantic import Field

from .base import ConfigModel
from .grpc import GRPCConfig
from .log import LogConfig


class AppConfig(ConfigModel):
    """Application configuration"""

    grpc: GRPCConfig = Field(default_factory=GRPCConfig)
    log: LogConfig = Field(default_factory=LogConfig)


def load_yaml_config(path: str) -> AppConfig:
    """Load a yaml config file and return a Config object"""
    return AppConfig.model_validate(_load_yaml(path))


def _load_yaml(path: str) -> dict[str, Any]:
    """Load a yaml config file and return a dict"""
    with open(path, "r") as f:
        config = yaml.safe_load(f)
    if not isinstance(config, dict):
        raise TypeError(f"Config file {path} has no top level mapping")
    return cast(dict[str, Any], config)
