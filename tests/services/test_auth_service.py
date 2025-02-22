from datetime import datetime, timedelta
from unittest.mock import Mock

import pytest

from src.config import AppConfig
from src.config.auth import AuthConfig, Credentials
from src.exception.base import Unauthorized
from src.models.session import Session
from src.repositories.session_repository import SessionRepository
from src.services.auth_service import AuthService


@pytest.fixture
def credentials() -> Credentials:
    return Credentials(username="test_user", password="test_pass")


@pytest.fixture
def auth_config(credentials: Credentials) -> AuthConfig:
    return AuthConfig(
        credentials=credentials,
        session_expiration_time=timedelta(seconds=10),
    )


@pytest.fixture
def app_config(auth_config: AuthConfig) -> AppConfig:
    return AppConfig(auth=auth_config)


@pytest.fixture
def session_repository() -> Mock:
    return Mock(spec=SessionRepository)


@pytest.fixture
def auth_service(app_config: AppConfig, session_repository: Mock) -> AuthService:
    return AuthService(app_config, session_repository)


def test_login_success(auth_service: AuthService, session_repository: Mock) -> None:
    session_repository.create_session.return_value = Mock(spec=Session)

    session = auth_service.login("test_user", "test_pass")

    assert session is not None
    session_repository.create_session.assert_called_once()


def test_login_invalid_credentials(auth_service: AuthService) -> None:
    with pytest.raises(Unauthorized, match="Invalid credentials"):
        auth_service.login("wrong_user", "wrong_pass")


def test_validate_session_success(
    auth_service: AuthService, session_repository: Mock
) -> None:
    now = datetime.now()
    session = Session(id="test_id", created_at=now, last_seen_at=now)
    session_repository.get_session.return_value = session

    auth_service.validate_session("test_id")

    session_repository.update_session.assert_called_once_with(session)


def test_validate_session_not_found(
    auth_service: AuthService, session_repository: Mock
) -> None:
    session_repository.get_session.return_value = None

    with pytest.raises(Unauthorized, match="Invalid session"):
        auth_service.validate_session("invalid_id")


def test_validate_session_expired(
    auth_service: AuthService, session_repository: Mock
) -> None:
    now = datetime.now()
    expired_time = now - timedelta(seconds=20)
    session = Session(id="test_id", created_at=expired_time, last_seen_at=expired_time)
    session_repository.get_session.return_value = session

    with pytest.raises(Unauthorized, match="Session expired"):
        auth_service.validate_session("test_id")
