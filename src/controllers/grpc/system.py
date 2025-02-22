import logging

import grpc  # type: ignore
import raybot.v1.system_pb2 as pb2
import raybot.v1.system_pb2_grpc as pb2_grpc  # type: ignore

logger = logging.getLogger(__name__)


class SystemService(pb2_grpc.SystemServiceServicer):
    def EmergencyStop(
        self,
        request: pb2.EmergencyStopRequest,
        context: grpc.ServicerContext,
    ) -> pb2.EmergencyStopResponse:
        logger.info(f"Emergency stop requested: {request}")
        return pb2.EmergencyStopResponse()
