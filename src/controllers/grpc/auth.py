import logging

import grpc  # type: ignore
import raybot.v1.auth_pb2 as pb2
import raybot.v1.auth_pb2_grpc as pb2_grpc  # type: ignore

logger = logging.getLogger(__name__)


class AuthService(pb2_grpc.AuthServiceServicer):
    def GetSession(
        self,
        request: pb2.GetSessionRequest,
        context: grpc.ServicerContext,
    ) -> pb2.GetSessionResponse:
        logger.info(f"Getting session: {request}")
        return pb2.GetSessionResponse()

    def register(self, srv: grpc.Server) -> None:
        pb2_grpc.add_AuthServiceServicer_to_server(self, srv)  # type: ignore
