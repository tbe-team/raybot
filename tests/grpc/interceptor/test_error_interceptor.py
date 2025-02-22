# type: ignore
from collections.abc import Generator
from concurrent import futures

import grpc
import pytest
import raybot.v1.auth_pb2 as pb2
import raybot.v1.auth_pb2_grpc as pb2_grpc

from src.controllers.grpc.interceptor.error_interceptor import ErrorInterceptor
from src.exception import BadRequest


class TestAuthService(pb2_grpc.AuthServiceServicer):
    def GetSession(
        self, request: pb2.GetSessionRequest, context: grpc.ServicerContext
    ) -> pb2.GetSessionResponse:
        if request.username == "error":
            raise BadRequest("test error")
        if request.username == "unhandled":
            raise ValueError("unhandled error")
        return pb2.GetSessionResponse()


@pytest.fixture(scope="module", autouse=True)
def grpc_server() -> Generator[str, None, None]:
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=1),
        interceptors=[ErrorInterceptor()],
    )
    pb2_grpc.add_AuthServiceServicer_to_server(TestAuthService(), server)
    port = server.add_insecure_port("[::]:0")

    server.start()

    yield f"localhost:{port}"

    server.stop(None)
    server.wait_for_termination()


@pytest.fixture(scope="module")
def server_client(grpc_server: str) -> Generator[pb2_grpc.AuthServiceStub, None, None]:
    channel = grpc.insecure_channel(grpc_server)
    stub = pb2_grpc.AuthServiceStub(channel)

    yield stub

    channel.close()


def test_handled_error(server_client: pb2_grpc.AuthServiceStub) -> None:
    with pytest.raises(grpc.RpcError) as exc:
        server_client.GetSession(pb2.GetSessionRequest(username="error"))

    assert exc.value.code() == grpc.StatusCode.INVALID_ARGUMENT
    assert "test error" in exc.value.details()


def test_unhandled_error(server_client: pb2_grpc.AuthServiceStub) -> None:
    with pytest.raises(grpc.RpcError) as exc:
        server_client.GetSession(pb2.GetSessionRequest(username="unhandled"))

    assert exc.value.code() == grpc.StatusCode.INTERNAL
    assert "Internal server error" in exc.value.details()


def test_no_error(server_client: pb2_grpc.AuthServiceStub) -> None:
    response = server_client.GetSession(pb2.GetSessionRequest(username="ok"))
    assert isinstance(response, pb2.GetSessionResponse)
