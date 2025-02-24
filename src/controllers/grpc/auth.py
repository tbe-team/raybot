import logging

import grpc  # type: ignore
import raybot.v1.auth_pb2 as pb2
import raybot.v1.auth_pb2_grpc as pb2_grpc  # type: ignore

from src.services.auth_service import AuthService as AuthServiceImpl

logger = logging.getLogger(__name__)


class AuthService(pb2_grpc.AuthServiceServicer):
    def __init__(self, auth_service: AuthServiceImpl) -> None:
        self._auth_service = auth_service
        super().__init__()

    def GetSession(
        self,
        request: pb2.GetSessionRequest,
        context: grpc.ServicerContext,
    ) -> pb2.GetSessionResponse:
        session = self._auth_service.login(request.username, request.password)
        return pb2.GetSessionResponse(
            session_id=session.id,
            heartbeat_window=int(
                self._auth_service.get_session_expiration_time().total_seconds()
            ),
        )

    def SendHeartbeat(
        self,
        request: pb2.SendHeartbeatRequest,
        context: grpc.ServicerContext,
    ) -> pb2.SendHeartbeatResponse:
        # We don't need to do anything here
        # Because session interceptor will take care of the session expiration
        return pb2.SendHeartbeatResponse()
