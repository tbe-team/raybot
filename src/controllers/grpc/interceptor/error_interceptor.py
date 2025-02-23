# type:ignore
from collections.abc import Callable
from typing import Any

import grpc
from grpc_interceptor import ServerInterceptor

from src.exception import RaybotException


class ErrorInterceptor(ServerInterceptor):
    """Error interceptor for gRPC server

    This interceptor catches RaybotException and returns the appropriate gRPC status code.
    If an exception is raised that is not a RaybotException, it returns an INTERNAL status code.

    """

    def intercept(
        self,
        method: Callable,
        request: Any,
        context: grpc.ServicerContext,
        method_name: str,
    ) -> Any:
        try:
            return method(request, context)
        except RaybotException as e:
            context.set_code(e.GRPC_STATUS)
            context.set_details(str(e))
            raise
        except Exception:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details("Internal server error")
            raise
