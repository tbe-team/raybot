from datetime import datetime

from src.config import AppConfig
from src.repositories.session_repository import SessionRepository


class BackgroundJobService:
    def __init__(
        self, config: AppConfig, session_repository: SessionRepository
    ) -> None:
        self._session_expiration_time = config.auth.session_expiration_time
        self._session_repo = session_repository

    def clean_up_expire_session(self) -> None:
        sessions = self._session_repo.get_all_session()
        now = datetime.now()
        expired = [
            s for s in sessions if s.last_seen_at < now - self._session_expiration_time
        ]
        for session in expired:
            self._session_repo.delete_session(session.id)
