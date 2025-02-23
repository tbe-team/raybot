# type: ignore
import logging
import threading
from concurrent import futures

import grpc
from raybot.v1 import (
    auth_pb2_grpc,
    drive_motor_pb2_grpc,
    lift_motor_pb2_grpc,
    state_pb2_grpc,
    system_pb2_grpc,
)

from src.app.application import Application
from src.controllers.grpc.auth import AuthService
from src.controllers.grpc.drive_motor import DriveMotorService
from src.controllers.grpc.interceptor import ErrorInterceptor, SessionInterceptor
from src.controllers.grpc.lift_motor import LiftMotorService
from src.controllers.grpc.robot_state import RobotStateService
from src.controllers.grpc.system import SystemService

logger = logging.getLogger(__name__)


def start_grpc_service(app: Application, interupt_event: threading.Event) -> None:
    """Start the grpc service.

    Blocking call. Should run as a separate thread.

    Args:
        config: The application config
        interupt_event: The event that will be set when the service should shutdown
    """
    srv = grpc.server(
        futures.ThreadPoolExecutor(max_workers=app.config.grpc.max_workers),
        interceptors=init_interceptors(app),
    )

    register_services(app, srv)

    grpc_port = app.config.grpc.port
    srv.add_insecure_port(f"[::]:{grpc_port}")

    logger.info("Starting grpc service at [::]:%d", grpc_port)
    srv.start()

    # Block until the interrupt event is set
    interupt_event.wait()

    logger.info("Shutting down grpc service...")
    srv.stop(None)
    logger.info("Grpc service shutdown complete")


def init_interceptors(app: Application) -> list[grpc.ServerInterceptor]:
    """Initialize the interceptors for the grpc service.

    Args:
        app: The application instance

    Returns:
        A list of interceptors
    """
    return [
        # ErrorInterceptor should be place first to catch all error.
        ErrorInterceptor(),
        SessionInterceptor(app.services.auth_service),
    ]


def register_services(app: Application, srv: grpc.Server) -> None:
    """Register all services to the grpc server.

    Args:
        app: The application instance
        srv: The grpc server
    """
    drive_motor_service = DriveMotorService()
    drive_motor_pb2_grpc.add_DriveMotorServiceServicer_to_server(
        drive_motor_service, srv
    )

    lift_motor_service = LiftMotorService()
    lift_motor_pb2_grpc.add_LiftMotorServiceServicer_to_server(lift_motor_service, srv)

    system_service = SystemService()
    system_pb2_grpc.add_SystemServiceServicer_to_server(system_service, srv)

    robot_state_service = RobotStateService()
    state_pb2_grpc.add_RobotStateServiceServicer_to_server(robot_state_service, srv)

    auth_service = AuthService(app.services.auth_service)
    auth_pb2_grpc.add_AuthServiceServicer_to_server(auth_service, srv)
