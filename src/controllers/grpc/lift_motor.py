import logging

import grpc  # type: ignore
import raybot.v1.lift_motor_pb2 as pb2
import raybot.v1.lift_motor_pb2_grpc as pb2_grpc  # type: ignore

logger = logging.getLogger(__name__)


class LiftMotorService(pb2_grpc.LiftMotorServiceServicer):
    def SetLiftMotorConfiguration(
        self,
        request: pb2.SetLiftMotorConfigurationRequest,
        context: grpc.ServicerContext,
    ) -> pb2.SetLiftMotorConfigurationResponse:
        logger.info(f"Setting lift motor configuration: {request}")
        return pb2.SetLiftMotorConfigurationResponse()
