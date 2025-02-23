# type:ignore
from collections.abc import Callable
from typing import Any, Protocol

import grpc
from grpc_interceptor import ServerInterceptor

from src.exception import Unauthorized

SESSION_METADATA_KEY = "session-id"
PUBLIC_METHODS = {"GetSession"}


class SessionValidator(Protocol):
    """Session validator interface"""

    def validate_session(self, session_id: str) -> None: ...


class SessionInterceptor(ServerInterceptor):
    """Session interceptor for gRPC server

    This interceptor validates the session ID for all methods that are not in the PUBLIC_METHODS set.

    """

    def __init__(self, session_validator: SessionValidator) -> None:
        self._session_validator = session_validator
        super().__init__()

    def intercept(
        self,
        method: Callable,
        request: Any,
        context: grpc.ServicerContext,
        method_name: str,
    ) -> Any:
        # Validate the last part of the method name
        method_suffix = method_name.split("/")[-1]
        if method_suffix in PUBLIC_METHODS:
            return method(request, context)

        # Get the session ID from the metadata
        metadata = dict(context.invocation_metadata())
        session_id = metadata.get(SESSION_METADATA_KEY)
        if session_id is None:
            raise Unauthorized("Session ID is required")

        # Validate the session
        self._session_validator.validate_session(session_id)

        return method(request, context)
