from datetime import datetime
from uuid import uuid4

from src.config import AppConfig
from src.exception import Unauthorized
from src.models.session import Session
from src.repositories.session_repository import SessionRepository


class AuthService:
    def __init__(
        self,
        config: AppConfig,
        session_repository: SessionRepository,
    ) -> None:
        self._credentials = config.auth.credentials
        self._session_expiration_time = config.auth.session_expiration_time
        self._session_repository = session_repository

    def login(self, username: str, password: str) -> Session:
        if (
            username != self._credentials.username
            or password != self._credentials.password
        ):
            raise Unauthorized("Invalid credentials")

        now = datetime.now()
        session = Session(
            id=str(uuid4()),
            created_at=now,
            last_seen_at=now,
        )
        return self._session_repository.create_session(session)

    def validate_session(self, session_id: str) -> None:
        session = self._session_repository.get_session(session_id)
        if session is None:
            raise Unauthorized("Invalid session")

        if session.last_seen_at < datetime.now() - self._session_expiration_time:
            raise Unauthorized("Session expired")

        session.last_seen_at = datetime.now()
        self._session_repository.update_session(session)
