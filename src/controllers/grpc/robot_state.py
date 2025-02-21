import logging

import grpc  # type: ignore
import raybot.v1.state_pb2 as pb2
import raybot.v1.state_pb2_grpc as pb2_grpc  # type: ignore

logger = logging.getLogger(__name__)


class RobotStateService(pb2_grpc.RobotStateServiceServicer):
    def GetRobotState(
        self, request: pb2.GetRobotStateRequest, context: grpc.ServicerContext
    ) -> pb2.GetRobotStateResponse:
        logger.info(f"Getting robot state: {request}")
        return pb2.GetRobotStateResponse()

    def register(self, srv: grpc.Server) -> None:
        pb2_grpc.add_RobotStateServiceServicer_to_server(self, srv)  # type: ignore
