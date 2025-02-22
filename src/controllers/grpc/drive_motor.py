import logging

import grpc  # type: ignore
import raybot.v1.drive_motor_pb2 as pb2
import raybot.v1.drive_motor_pb2_grpc as pb2_grpc  # type: ignore

logger = logging.getLogger(__name__)


class DriveMotorService(pb2_grpc.DriveMotorServiceServicer):
    def SetDriveMotorConfiguration(
        self,
        request: pb2.SetDriveMotorConfigurationRequest,
        context: grpc.ServicerContext,
    ) -> pb2.SetDriveMotorConfigurationResponse:
        logger.info(f"Setting drive motor configuration: {request}")
        return pb2.SetDriveMotorConfigurationResponse()
