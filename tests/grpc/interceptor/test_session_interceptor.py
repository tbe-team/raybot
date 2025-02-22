# type: ignore
from collections.abc import Generator
from concurrent import futures
from unittest.mock import Mock

import grpc
import pytest
from raybot.v1 import auth_pb2, auth_pb2_grpc, system_pb2, system_pb2_grpc

from src.controllers.grpc.interceptor import ErrorInterceptor, SessionInterceptor
from src.controllers.grpc.interceptor.session_interceptor import PUBLIC_METHODS
from src.services.session import SessionService


class TestPublicService(auth_pb2_grpc.AuthServiceServicer):
    def GetSession(
        self, request: auth_pb2.GetSessionRequest, context: grpc.ServicerContext
    ) -> auth_pb2.GetSessionResponse:
        return auth_pb2.GetSessionResponse(session_id="test")


class TestPrivateService(system_pb2_grpc.SystemServiceServicer):
    def EmergencyStop(
        self, request: system_pb2.EmergencyStopRequest, context: grpc.ServicerContext
    ) -> system_pb2.EmergencyStopResponse:
        return system_pb2.EmergencyStopResponse()


@pytest.fixture(scope="module")
def session_service() -> Mock:
    return Mock(spec=SessionService)


@pytest.fixture(autouse=True)
def reset_mocks(session_service: Mock) -> None:
    session_service.reset_mock()


@pytest.fixture(scope="module")
def grpc_server(session_service: SessionService) -> Generator[str, None, None]:
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=1),
        interceptors=[
            ErrorInterceptor(),  # We need to catch all exceptions
            SessionInterceptor(session_service),
        ],
    )
    auth_pb2_grpc.add_AuthServiceServicer_to_server(TestPublicService(), server)
    system_pb2_grpc.add_SystemServiceServicer_to_server(TestPrivateService(), server)
    port = server.add_insecure_port("[::]:0")

    server.start()

    yield f"localhost:{port}"

    server.stop(None)
    server.wait_for_termination()


@pytest.fixture
def public_client(
    grpc_server: str,
) -> Generator[auth_pb2_grpc.AuthServiceStub, None, None]:
    channel = grpc.insecure_channel(grpc_server)
    stub = auth_pb2_grpc.AuthServiceStub(channel)

    yield stub

    channel.close()


@pytest.fixture
def private_client(
    grpc_server: str,
) -> Generator[system_pb2_grpc.SystemServiceStub, None, None]:
    channel = grpc.insecure_channel(grpc_server)
    stub = system_pb2_grpc.SystemServiceStub(channel)

    yield stub

    channel.close()


def test_public_method() -> None:
    assert len(PUBLIC_METHODS) == 1
    assert "GetSession" in PUBLIC_METHODS


def test_public_service_not_raise_exception(
    public_client: auth_pb2_grpc.AuthServiceStub,
) -> None:
    response = public_client.GetSession(auth_pb2.GetSessionRequest(username="test"))
    assert response.session_id == "test"


def test_private_service_with_session(
    private_client: system_pb2_grpc.SystemServiceStub,
    session_service: Mock,
) -> None:
    session_service.validate_session.return_value = None

    metadata = (("session-id", "test"),)
    private_client.EmergencyStop(system_pb2.EmergencyStopRequest(), metadata=metadata)

    session_service.validate_session.assert_called_once()


def test_private_service_without_session(
    private_client: system_pb2_grpc.SystemServiceStub,
    session_service: Mock,
) -> None:
    with pytest.raises(grpc.RpcError) as exc:
        private_client.EmergencyStop(system_pb2.EmergencyStopRequest())

    assert exc.value.code() == grpc.StatusCode.UNAUTHENTICATED
    assert "Session ID is required" in exc.value.details()

    session_service.validate_session.assert_not_called()


def test_private_service_with_invalid_session_key(
    private_client: system_pb2_grpc.SystemServiceStub,
    session_service: Mock,
) -> None:
    with pytest.raises(grpc.RpcError) as exc:
        metadata = (("invalid-session-key", "invalid"),)
        private_client.EmergencyStop(
            system_pb2.EmergencyStopRequest(), metadata=metadata
        )

    assert exc.value.code() == grpc.StatusCode.UNAUTHENTICATED
    assert "Session ID is required" in exc.value.details()

    session_service.validate_session.assert_not_called()
