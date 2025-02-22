from datetime import datetime, timedelta
from unittest.mock import Mock

import pytest

from src.config import AppConfig, AuthConfig
from src.config.auth import Credentials
from src.models.session import Session
from src.repositories.session_repository import SessionRepository
from src.services.background_job_service import BackgroundJobService


@pytest.fixture
def auth_config() -> AuthConfig:
    return AuthConfig(
        credentials=Mock(spec=Credentials),
        session_expiration_time=timedelta(seconds=10),
    )


@pytest.fixture
def app_config(auth_config: AuthConfig) -> AppConfig:
    return AppConfig(auth=auth_config)


@pytest.fixture
def session_repository() -> Mock:
    return Mock(spec=SessionRepository)


@pytest.fixture
def background_job_service(
    app_config: AppConfig, session_repository: Mock
) -> BackgroundJobService:
    return BackgroundJobService(app_config, session_repository)


def test_clean_up_expire_session(
    background_job_service: BackgroundJobService, session_repository: Mock
) -> None:
    now = datetime.now()
    active_session = Session(id="active", created_at=now, last_seen_at=now)
    expired_session = Session(
        id="expired",
        created_at=now - timedelta(seconds=20),
        last_seen_at=now - timedelta(seconds=20),
    )

    session_repository.get_all_session.return_value = [active_session, expired_session]

    background_job_service.clean_up_expire_session()

    session_repository.delete_session.assert_called_once_with("expired")


def test_clean_up_no_expired_sessions(
    background_job_service: BackgroundJobService, session_repository: Mock
) -> None:
    now = datetime.now()
    active_session = Session(id="active", created_at=now, last_seen_at=now)

    session_repository.get_all_session.return_value = [active_session]

    background_job_service.clean_up_expire_session()

    session_repository.delete_session.assert_not_called()


def test_clean_up_empty_sessions(
    background_job_service: BackgroundJobService, session_repository: Mock
) -> None:
    session_repository.get_all_session.return_value = []

    background_job_service.clean_up_expire_session()

    session_repository.delete_session.assert_not_called()
