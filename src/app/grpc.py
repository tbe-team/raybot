import logging
import threading
from concurrent import futures

import grpc  # type: ignore

from src.config.grpc import GRPCConfig
from src.controllers.grpc.auth import AuthService
from src.controllers.grpc.drive_motor import DriveMotorService
from src.controllers.grpc.lift_motor import LiftMotorService
from src.controllers.grpc.robot_state import RobotStateService
from src.controllers.grpc.system import SystemService

logger = logging.getLogger(__name__)


def start_grpc_service(config: GRPCConfig, interupt_event: threading.Event) -> None:
    """Start the grpc service.

    Blocking call. Should run as a separate thread.

    Args:
        config: The application config
        interupt_event: The event that will be set when the service should shutdown
    """
    srv = grpc.server(futures.ThreadPoolExecutor(max_workers=10))  # type: ignore

    drive_motor_service = DriveMotorService()
    drive_motor_service.register(srv)

    lift_motor_service = LiftMotorService()
    lift_motor_service.register(srv)

    system_service = SystemService()
    system_service.register(srv)

    robot_state_service = RobotStateService()
    robot_state_service.register(srv)

    auth_service = AuthService()
    auth_service.register(srv)

    srv.add_insecure_port(f"[::]:{config.port}")

    logger.info("Starting grpc server at [::]:%d", config.port)
    srv.start()

    # Block until the interrupt event is set
    interupt_event.wait()

    logger.info("Shutting down grpc server...")
    srv.stop(None)
    logger.info("Grpc server shutdown complete")
