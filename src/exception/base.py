import grpc  # type: ignore


class RaybotException(Exception):
    """Base exception for all Raybot exceptions"""

    GRPC_STATUS: grpc.StatusCode = grpc.StatusCode.UNKNOWN

    def __init__(self, message: str = "Unknown exception") -> None:
        super().__init__(message)


class InternalServerError(RaybotException):
    """Internal server error"""

    pass


class Unauthorized(RaybotException):
    """Unauthorized"""

    GRPC_STATUS = grpc.StatusCode.UNAUTHENTICATED

    def __init__(self, message: str = "Unauthorized") -> None:
        super().__init__(message)


class PermissionDenied(RaybotException):
    """Permission denied"""

    GRPC_STATUS = grpc.StatusCode.PERMISSION_DENIED

    def __init__(self, message: str = "Permission denied") -> None:
        super().__init__(message)


class BadRequest(RaybotException):
    """Bad request"""

    GRPC_STATUS = grpc.StatusCode.INVALID_ARGUMENT

    def __init__(self, message: str = "Bad request") -> None:
        super().__init__(message)


class NotFound(RaybotException):
    """Not found"""

    GRPC_STATUS = grpc.StatusCode.NOT_FOUND

    def __init__(self, message: str = "Not found") -> None:
        super().__init__(message)


class AlreadyExists(RaybotException):
    """Already exists"""

    GRPC_STATUS = grpc.StatusCode.ALREADY_EXISTS

    def __init__(self, message: str = "Already exists") -> None:
        super().__init__(message)
